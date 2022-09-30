package impl

import (
	"context"
	"strconv"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/core/module/trace"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

type ServiceImpl struct {
}

func NewServiceImpl() *ServiceImpl {
	r := &ServiceImpl{}
	return r
}

func (s *ServiceImpl) ReportOnlineMessage(ctx context.Context) {
	var log *logger.Logger
	if ctx == nil {
		log = logger.New("ReportOnlineMessage Start")
	} else {
		log = logger.ReqLogger(ctx)
	}
	db := mysql.GetLive(log.ReqID())
	var numLives int
	db.Table("live_entities").Count(&numLives)
	var numUsers int
	db.Table("live_users").Count(&numUsers)
	var value = map[string]string{
		"lives": strconv.Itoa(numLives),
		"users": strconv.Itoa(numUsers),
	}
	trace.ReportEvent(ctx, "overview", value)
}

type CommonSingleStats struct {
	Type            int    `json:"type"`
	TypeDescription string `json:"type_description"`
	PageView        int    `json:"page_view"`
	UniqueVisitor   int    `json:"unique_visitor"`
}

type CommonStats struct {
	Flow string              `json:"flow"`
	Info []CommonSingleStats `json:"info"`
}

func (s *ServiceImpl) GetStatsSingleLive(ctx context.Context, liveId string) (*CommonStats, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	info := make([]CommonSingleStats, len(model.StatsTypeDescription))
	data := &CommonStats{
		Info: info,
	}
	for i := 0; i < len(model.StatsTypeDescription); i++ {
		data.Info[i].Type = i + 1
		data.Info[i].TypeDescription = model.StatsTypeDescription[i+1]
		var pv *int
		var uv *int
		db.DB().QueryRow("SELECT count(DISTINCT(user_id)) FROM stats_single_live  WHERE type = ? and live_id = ?", i+1, liveId).Scan(&uv)
		sql := "select sum(count) from stats_single_live where type = ? and live_id = ? ;"
		err := db.DB().QueryRow(sql, i+1, liveId).Scan(&pv)
		if err != nil {
			return nil, err
		}
		if pv != nil {
			data.Info[i].PageView = *pv
		}
		if uv != nil {
			data.Info[i].UniqueVisitor = *uv
		}
	}
	return data, nil
}

func (s *ServiceImpl) PostStatsSingleLive(ctx context.Context, entities []*model.StatsSingleLiveEntity) error {
	var err error
	for _, entity := range entities {
		err = s.UpdateSingleLive(ctx, entity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ServiceImpl) UpdateSingleLive(ctx context.Context, entity *model.StatsSingleLiveEntity) error {
	log := logger.ReqLogger(ctx)

	db := mysql.GetLive(log.ReqID())
	db = db.Model(model.StatsSingleLiveEntity{}).Where("live_id = ? and biz_id = ? and user_id = ? and type = ?", entity.LiveId, entity.BizId, entity.UserId, entity.Type)

	old := model.StatsSingleLiveEntity{}
	result := db.First(&old)
	if result.Error != nil {
		if result.RecordNotFound() {
			return s.createSingleLiveInfo(ctx, entity)
		} else {
			log.Errorf("find old stats single live error %+v", result.Error)
			return api.ErrDatabase
		}
	}

	updates := singleLiveUpdates(&old, entity)
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = timestamp.Now()
	result = db.Update(updates)
	if result.Error != nil {
		log.Errorf("update stats single live error %v", result.Error)
		return api.ErrDatabase
	} else {
		return nil
	}
}

func singleLiveUpdates(old, update *model.StatsSingleLiveEntity) map[string]interface{} {
	updates := map[string]interface{}{}
	if update.Count > 0 {
		updates["count"] = old.Count + update.Count
	}
	return updates
}

func (s *ServiceImpl) createSingleLiveInfo(ctx context.Context, entity *model.StatsSingleLiveEntity) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	result := db.Create(entity)
	if result.Error != nil {
		log.Errorf("create single live info %+v, error %+v", entity, result.Error)
		return api.ErrDatabase
	}
	return nil
}

func (s *ServiceImpl) SaveStatsSingleLive(ctx context.Context, entities []*model.StatsSingleLiveEntity) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	tx := db.Begin()

	for _, entity := range entities {
		sql, params := entity.SaveSql()
		if err := tx.Exec(sql, params...).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
