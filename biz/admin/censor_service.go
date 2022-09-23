package admin

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

type CCService interface {
	UpdateCensorConfig(ctx context.Context, mod *model.CensorConfig) error
	GetCensorConfig(ctx context.Context) (*model.CensorConfig, error)
	CreateCensorJob(ctx context.Context, liveEntity *model.LiveEntity) error
	StopCensorJob(ctx context.Context, liveId string) error
	GetLiveCensorJobByLiveId(ctx context.Context, liveId string) (*model.LiveCensor, error)
	GetLiveCensorJobByJobId(ctx context.Context, jobId string) (*model.LiveCensor, error)

	GetCensorImageById(ctx context.Context, imageId uint) (*model.CensorImage, error)
	SaveLiveCensorJob(ctx context.Context, liveId string, jobId string, config *model.CensorConfig) error
	SaveCensorImage(ctx context.Context, image *model.CensorImage) error
	BatchUpdateCensorImage(ctx context.Context, images []uint, updates map[string]interface{}) error
	SearchCensorImage(ctx context.Context, isReview, pageNum, pageSize int, liveId string) (image []model.CensorImage, totalCount int, err error)
	SearchCensorLive(ctx context.Context, audit, pageNum, pageSize int) (censorLive []CensorLive, totalCount int, err error)
	GetUnauditCount(ctx context.Context, liveId string) (len int, err error)
}

type CensorService struct {
}

var cService CCService = &CensorService{}

func GetCensorService() CCService {
	return cService
}

func (c *CensorService) UpdateCensorConfig(ctx context.Context, mod *model.CensorConfig) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	mod.ID = 1
	err := db.Model(model.CensorConfig{}).Save(mod).Error
	if err != nil {
		return api.ErrDatabase
	}
	return nil
}

func (c *CensorService) GetCensorConfig(ctx context.Context) (*model.CensorConfig, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	mod := &model.CensorConfig{}
	result := db.Model(model.CensorConfig{}).First(mod)
	if result.Error != nil {
		if result.RecordNotFound() {
			mod.Enable = true
			mod.ID = 1
			mod.Pulp = true
			mod.Interval = 20
			err := db.Model(model.CensorConfig{}).Save(mod).Error
			if err != nil {
				return nil, api.ErrDatabase
			}
		} else {
			return nil, api.ErrDatabase
		}
	}
	return mod, nil
}

func (c *CensorService) CreateCensorJob(ctx context.Context, liveEntity *model.LiveEntity) error {
	log := logger.ReqLogger(ctx)
	config, err := c.GetCensorConfig(ctx)
	if err != nil {
		log.Errorf("GetCensorConfig Error %v", err)
		return err
	}
	if config.Enable == false {
		return nil
	}
	resp, err := GetJobService().JobCreate(ctx, liveEntity, config)
	if err != nil {
		log.Errorf("JobCreate Error %v", err)
		return err
	}
	err = c.SaveLiveCensorJob(ctx, liveEntity.LiveId, resp.Data.JobID, config)
	if err != nil {
		log.Errorf("SaveLiveCensorJob Error %v", err)
		return nil
	}
	return nil
}

func (c *CensorService) SaveLiveCensorJob(ctx context.Context, liveId string, jobId string, config *model.CensorConfig) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	m := &model.LiveCensor{
		LiveID:     liveId,
		JobID:      jobId,
		Interval:   config.Interval,
		Politician: config.Politician,
		Pulp:       config.Pulp,
		Ads:        config.Ads,
		Terror:     config.Terror,
	}
	err := db.Model(model.LiveCensor{}).Save(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CensorService) GetLiveCensorJobByLiveId(ctx context.Context, liveId string) (*model.LiveCensor, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	m := &model.LiveCensor{}
	err := db.Model(model.LiveCensor{}).First(m, "live_id = ?", liveId).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *CensorService) StopCensorJob(ctx context.Context, liveId string) error {
	log := logger.ReqLogger(ctx)
	liveCensorJob, err := c.GetLiveCensorJobByLiveId(ctx, liveId)
	if err != nil {
		log.Errorf("GetLiveCensorJobByLiveId Error %v", err)
		return err
	}

	req := &JobCreateResponseData{
		JobID: liveCensorJob.JobID,
	}
	err = GetJobService().JobClose(ctx, req)
	if err != nil {
		log.Errorf("post StopCensorJob Error %v", err)
		return err
	}
	return nil
}

func (c *CensorService) GetCensorImageById(ctx context.Context, imageId uint) (*model.CensorImage, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	m := &model.CensorImage{}
	err := db.Model(model.CensorImage{}).First(m, "id = ?", imageId).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *CensorService) GetLiveCensorJobByJobId(ctx context.Context, jobId string) (*model.LiveCensor, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	m := &model.LiveCensor{}
	err := db.Model(model.LiveCensor{}).First(m, "job_id = ?", jobId).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *CensorService) SaveCensorImage(ctx context.Context, image *model.CensorImage) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	err := db.Model(model.CensorImage{}).Save(image).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CensorService) GetUnauditCount(ctx context.Context, liveId string) (len int, err error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLiveReadOnly(log.ReqID())
	err = db.Model(&model.CensorImage{}).Where(" is_review = ? and live_id = ? ", 0, liveId).Count(&len).Error
	return
}

// SearchCensorImage 0： 没审核 1：审核 2：都需要list出来/*
func (c *CensorService) SearchCensorImage(ctx context.Context, isReview, pageNum, pageSize int, liveId string) (image []model.CensorImage, totalCount int, err error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLiveReadOnly(log.ReqID())
	image = make([]model.CensorImage, 0)

	var where *gorm.DB
	var all *gorm.DB
	if liveId == "" {
		where = db.Model(&model.CensorImage{}).Where(" is_review = ? ", isReview)
		all = db.Model(&model.CensorImage{})
	} else {
		where = db.Model(&model.CensorImage{}).Where(" is_review = ? and live_id = ?", isReview, liveId)
		all = db.Model(&model.CensorImage{}).Where(" live_id = ? ", liveId)
	}

	if isReview == 0 {
		err = where.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&image).Error
		err = where.Count(&totalCount).Error
	} else if isReview == 1 {
		err = where.Count(&totalCount).Error
		err = where.Order("review_answer desc").Order("created_at desc ").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&image).Error
	} else {
		err = all.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&image).Error
		err = all.Count(&totalCount).Error
	}

	if err != nil {
		log.Errorf("search CensorImage error: %v", err)
		return
	}
	return
}

const AuditNo = 1
const AuditAll = 0

func (c *CensorService) SearchCensorLive(ctx context.Context, audit, pageNum, pageSize int) (censorLive []CensorLive, totalCount int, err error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLiveReadOnly(log.ReqID())
	lives := make([]model.LiveEntity, 0)
	var db2 *gorm.DB
	if audit == AuditNo {
		db2 = db.Model(&model.LiveEntity{}).Where("unaudit_censor_count > 0").Where("stop_reason != ?", model.LiveStopReasonCensor)
	} else {
		db2 = db.Model(&model.LiveEntity{}).Where("unaudit_censor_count >= 0")
	}
	err = db2.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&lives).Error
	err = db2.Count(&totalCount).Error
	if err != nil {
		log.Errorf("SearchCensorLive %v", err)
		return nil, 0, err
	}

	for _, live := range lives {
		cl := CensorLive{
			LiveId:     live.LiveId,
			Title:      live.Title,
			AnchorId:   live.AnchorId,
			Status:     live.Status,
			Count:      live.UnauditCensorCount,
			Time:       live.LastCensorTime,
			StopReason: live.StopReason,
			StopAt:     live.StopAt,
		}
		censorLive = append(censorLive, cl)
	}
	return
}

func (c *CensorService) BatchUpdateCensorImage(ctx context.Context, images []uint, updates map[string]interface{}) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	db = db.Model(model.CensorImage{})

	result := db.Where(" id in (?) ", images).Update(updates)
	if result.Error != nil {
		log.Errorf("update user error %v", result.Error)
		return api.ErrDatabase
	} else {
		return nil
	}
}

type CensorLive struct {
	LiveId       string               `json:"live_id"`
	Title        string               `json:"title"`
	AnchorId     string               `json:"anchor_id"`
	Nick         string               `json:"nick"`
	Status       int                  `json:"live_status"`
	AnchorStatus int                  `json:"anchor_status"`
	StopReason   string               `json:"stop_reason"`
	StopAt       *timestamp.Timestamp `json:"stop_at"`
	Count        int                  `json:"count"`
	Time         timestamp.Timestamp  `json:"time"`
}

type JobQueryRequest struct {
	Job         string   `json:"job"`
	Suggestions []string `json:"suggestions"`
	Start       int64    `json:"start"`
	End         int64    `json:"end"`
}

type JobQueryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Marker string `json:"marker"`
		Items  struct {
			Image []struct {
				Code      int    `json:"code,omitempty"`
				Message   string `json:"message,omitempty"`
				Job       string `json:"job,omitempty"`
				Timestamp int    `json:"timestamp,omitempty"`
				Url       string `json:"url,omitempty"`
				Result    struct {
					Suggestion string `json:"suggestion"`
					Scenes     struct {
						Pulp struct {
							Suggestion string `json:"suggestion"`
							Details    []struct {
								Suggestion string  `json:"suggestion,omitempty"`
								Label      string  `json:"label,omitempty"`
								Score      float64 `json:"score,omitempty"`
							} `json:"details"`
						} `json:"pulp"`
					} `json:"scenes"`
				} `json:"result,omitempty"`
			} `json:"image"`
			Audio []struct {
				Code      int    `json:"code,omitempty"`
				Message   string `json:"message,omitempty"`
				Job       string `json:"job,omitempty"`
				Start     int    `json:"start,omitempty"`
				End       int    `json:"end,omitempty"`
				Url       string `json:"url,omitempty"`
				AudioText string `json:"audio_text,omitempty"`
				Result    struct {
					Suggestion string `json:"suggestion"`
					Scenes     struct {
						Antispam struct {
							Suggestion string `json:"suggestion"`
							Details    []struct {
								Suggestion string  `json:"suggestion"`
								Label      string  `json:"label"`
								Text       string  `json:"text"`
								Score      float64 `json:"score"`
							} `json:"details"`
						} `json:"antispam"`
					} `json:"scenes"`
				} `json:"result,omitempty"`
			} `json:"audio"`
		} `json:"items"`
	} `json:"data"`
}

type JobListRequest struct {
	Start  int64  `json:"start"`
	End    int64  `json:"end"`
	Status string `json:"status"`
	Limit  int    `json:"limit"`
	Marker string `json:"marker"`
}

type JobListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Marker string `json:"marker"`
		Items  []struct {
			Id   string `json:"id"`
			Data struct {
				Id   string `json:"id"`
				Uri  string `json:"uri"`
				Info string `json:"info"`
			} `json:"data"`
			Params struct {
				HookUrl  string `json:"hook_url"`
				HookAuth bool   `json:"hook_auth"`
				Image    struct {
					IsOn          bool     `json:"is_on"`
					Scenes        []string `json:"scenes"`
					IntervalMsecs int      `json:"interval_msecs"`
					Saver         struct {
						Uid    int    `json:"uid"`
						Bucket string `json:"bucket"`
						Prefix string `json:"prefix"`
					} `json:"saver"`
					HookRule int `json:"hook_rule"`
				} `json:"image"`
			} `json:"params"`
			Message   string `json:"message"`
			Status    string `json:"status"`
			CreatedAt int    `json:"created_at"`
			UpdatedAt int    `json:"updated_at"`
		} `json:"items"`
	} `json:"data"`
}

type JobCreateRequest struct {
	Data   JobLiveData     `json:"data"`
	Params JobCreateParams `json:"params"`
}

type JobLiveData struct {
	ID   string `json:"ID"`
	Url  string `json:"uri"`
	Info string `json:"info"`
}

type JobCreateParams struct {
	Image    JobImage `json:"image"`
	HookUrl  string   `json:"hook_url"`
	HookAuth bool     `json:"hook_auth"`
}

type JobImage struct {
	IsOn          bool          `json:"is_on"`
	Scenes        []string      `json:"scenes"`
	IntervalMsecs int           `json:"interval_msecs"`
	Saver         JobImageSaver `json:"saver"`
	HookRule      int           `json:"hook_rule"`
}

type JobImageSaver struct {
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
}

type JobCreateResponse struct {
	RequestId string                `json:"request_id"` //请求ID
	Code      int                   `json:"code"`       //错误码，0 成功，其他失败
	Message   string                `json:"message"`    //错误信息
	Data      JobCreateResponseData `json:"data"`
}

type JobCreateResponseData struct {
	JobID string `json:"job"`
}
