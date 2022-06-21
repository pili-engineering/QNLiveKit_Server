package live

import (
	"context"
	"errors"
	"time"

	"github.com/qbox/livekit/common/api"

	"github.com/qbox/livekit/common/im"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/user"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/common/rtc"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
	"github.com/qbox/livekit/utils/uuid"
)

type IService interface {
	CreateLive(context context.Context, req *CreateLiveRequest, anchorId string) (live *model.LiveEntity, err error)

	DeleteLive(context context.Context, liveId string, anchorId string) (err error)

	StartLive(context context.Context, liveId string, anchorId string) (roomToken string, err error)

	StopLive(context context.Context, liveId string, anchorId string) (err error)

	LiveInfo(context context.Context, liveId string) (live *model.LiveEntity, err error)

	LiveList(context context.Context, pageNum, pageSize int) (lives []model.LiveEntity, totalCount int, err error)

	LiveUserList(context context.Context, liveId string, pageNum, pageSize int) (users []model.LiveRoomUserEntity, totalCount int, err error)

	UpdateExtends(context context.Context, liveId string, extends model.Extends) (err error)

	JoinLiveRoom(context context.Context, liveId string, userId string) (err error)

	LeaveLiveRoom(context context.Context, liveId string, userId string) (err error)

	SearchLive(context context.Context, keyword string, flag, pageNum, pageSize int) (lives []model.LiveEntity, totalCount int, err error)

	// CurrentLiveRoom 查找主播当前所在的直播间
	CurrentLiveRoom(ctx context.Context, userId string) (liveEntity *model.LiveEntity, err error)

	Heartbeat(context context.Context, liveId string, userId string) (liveEntity *model.LiveEntity, err error)

	// StartRelay 绑定直播间与跨房PK 会话
	StartRelay(ctx context.Context, roomId, userId string, sid string) (err error)

	// StopRelay 解绑指定的跨房PK 会话
	StopRelay(ctx context.Context, roomId, userId string, sid string) (err error)

	TimeoutLiveUser(ctx context.Context, now time.Time)

	TimeoutLiveRoom(ctx context.Context, now time.Time)

	FindLiveRoomUser(context context.Context, liveId string, userId string) (liveRoomUser *model.LiveRoomUserEntity, err error)
}

type Service struct {
}

var service IService = &Service{}

func GetService() IService {
	return service
}

type CreateLiveRequest struct {
	Title    string        `json:"title"`
	Notice   string        `json:"notice"`
	CoverUrl string        `json:"cover_url"`
	Extends  model.Extends `json:"extends"`
}

func (s *Service) CreateLive(context context.Context, req *CreateLiveRequest, anchorId string) (live *model.LiveEntity, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	liveId := uuid.Gen()

	liveUser, err := user.GetService().FindUser(context, anchorId)
	if err != nil {
		log.Errorf("create live failed, user not found, userId: %s, err: %v", anchorId, err)
		return
	}
	rtcClient := rtc.GetService()
	imClient := im.GetService()
	chatroom, err := imClient.CreateChatroom(context, liveUser.ImUserid, liveId)
	if err != nil {
		log.Errorf("create chatroom failed, err: %v", err)
		return
	}
	live = &model.LiveEntity{
		LiveId:      liveId,
		Title:       req.Title,
		Notice:      req.Notice,
		CoverUrl:    req.CoverUrl,
		Extends:     req.Extends,
		AnchorId:    anchorId,
		Status:      model.LiveStatusPrepare,
		PkId:        "",
		OnlineCount: 0,
		StartAt:     timestamp.Now(),
		EndAt:       timestamp.Now(),
		ChatId:      chatroom,
		PushUrl:     rtcClient.StreamPubURL(liveId),
		RtmpPlayUrl: rtcClient.StreamRtmpPlayURL(liveId),
		FlvPlayUrl:  rtcClient.StreamFlvPlayURL(liveId),
		HlsPlayUrl:  rtcClient.StreamHlsPlayURL(liveId),
	}
	err = db.Create(live).Error
	return
}

func (s *Service) DeleteLive(context context.Context, liveId string, anchorId string) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err = db.Delete(&model.LiveEntity{}, "live_id = ? and anchor_id = ? and status = ?", liveId, anchorId, model.LiveStatusOn).Error
	return
}

func (s *Service) StartLive(context context.Context, liveId string, anchorId string) (roomToken string, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	live, err := s.LiveInfo(context, liveId)

	rtcClient := rtc.GetService()
	if err != nil {
		log.Errorf("LiveInfo error:%v", err)
		return
	}
	if live.Status != model.LiveStatusPrepare {
		err = errors.New("live status error")
		return
	}
	if live.AnchorId != anchorId {
		err = errors.New("user not anchor")
		return
	}

	////判断主播不在其他直播间
	liveUser, err := s.getOrCreateLiveRoomUser(context, anchorId)
	if err != nil {
		return "", err
	}

	live.Status = model.LiveStatusOn
	live.StartAt = timestamp.Now()
	live.LastHeartbeatAt = timestamp.Now()
	live.UpdatedAt = timestamp.Now()
	err = db.Save(live).Error
	roomToken = rtcClient.GetRoomToken(anchorId, liveId)

	now := timestamp.Now()
	liveUser.Status = model.LiveRoomUserStatusOnline
	liveUser.LiveId = liveId
	liveUser.UpdatedAt = now
	liveUser.HeartBeatAt = &now

	db.Save(liveUser)

	return
}

func (s *Service) StopLive(context context.Context, liveId string, anchorId string) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	live, err := s.LiveInfo(context, liveId)
	if err != nil {
		log.Errorf("LiveInfo error:%v", err)
		return
	}
	if live.Status != model.LiveStatusOn {
		err = errors.New("live status error")
		return
	}
	if live.AnchorId != anchorId {
		err = errors.New("user not anchor")
		return
	}
	live.Status = model.LiveStatusOff
	live.EndAt = timestamp.Now()
	err = db.Save(live).Error
	return
}

func (s *Service) LiveInfo(context context.Context, liveId string) (live *model.LiveEntity, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	live = &model.LiveEntity{}
	err = db.Where("live_id = ? ", liveId).First(live).Error
	return
}

func (s *Service) LiveList(context context.Context, pageNum, pageSize int) (lives []model.LiveEntity, totalCount int, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	lives = make([]model.LiveEntity, 0)
	err = db.Where("status = ?", model.LiveStatusOn).Order("updated_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&lives).Error
	err = db.Model(&model.LiveEntity{}).Where("status = ?", model.LiveStatusOn).Count(&totalCount).Error
	return
}

func (s *Service) LiveUserList(context context.Context, liveId string, pageNum, pageSize int) (users []model.LiveRoomUserEntity, totalCount int, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	users = make([]model.LiveRoomUserEntity, 0)
	err = db.Where("live_id = ?  and status = ?", liveId, model.LiveRoomUserStatusOnline).Order("updated_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&users).Error
	err = db.Model(&model.LiveRoomUserEntity{}).Where("live_id = ?  and status = ?", liveId, model.LiveRoomUserStatusOnline).Count(&totalCount).Error
	return
}

func (s *Service) UpdateExtends(context context.Context, liveId string, extends model.Extends) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	live, err := s.LiveInfo(context, liveId)
	if err != nil {
		log.Errorf("get live error: %v", err)
		return
	}
	live.Extends = extends
	err = db.Save(live).Error
	return
}

func (s *Service) JoinLiveRoom(context context.Context, liveId string, userId string) (err error) {
	log := logger.ReqLogger(context)

	liveUser, err := s.getOrCreateLiveRoomUser(context, userId)
	if err != nil {
		log.Errorf("get live room user error")
		return err
	}

	if liveUser.Status == model.LiveRoomUserStatusOnline && liveUser.LiveId != liveId {
		log.Infof("first")
		s.LeaveLiveRoom(context, liveUser.LiveId, userId)
	}

	now := timestamp.Now()
	liveUser.Status = model.LiveRoomUserStatusOnline
	liveUser.LiveId = liveId
	liveUser.HeartBeatAt = &now
	liveUser.UpdatedAt = timestamp.Now()

	db := mysql.GetLive(log.ReqID())
	err = db.Save(liveUser).Error
	if err != nil {
		log.Errorf("save live user error %v", err)
	}

	return err
}

func (s *Service) LeaveLiveRoom(context context.Context, liveId string, userId string) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	liveRoomUser := &model.LiveRoomUserEntity{}
	result := db.Where("live_id = ? and user_id = ? ", liveId, userId).First(liveRoomUser)
	if result.Error != nil {
		err = result.Error
	} else {
		liveRoomUser.LiveId = ""
		liveRoomUser.Status = model.LiveRoomUserStatusLeave
		liveRoomUser.UpdatedAt = timestamp.Now()
		err = db.Save(liveRoomUser).Error
	}
	return
}

func (s *Service) FindLiveRoomUser(context context.Context, liveId string, userId string) (liveRoomUser *model.LiveRoomUserEntity, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	liveRoomUser = &model.LiveRoomUserEntity{}
	result := db.Where("live_id = ? and user_id = ?", liveId, userId).First(liveRoomUser)
	if result.Error != nil {
		if result.RecordNotFound() {
			err = api.ErrNotFound
		} else {
			log.Errorf("find live room user error %v", err)
			err = api.ErrDatabase
		}
	}
	return
}

func (s *Service) FindUserLive(context context.Context, liveId string, userInfo *liveauth.UserInfo) (liveRoomUser *model.LiveRoomUserEntity, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	liveRoomUser = &model.LiveRoomUserEntity{}
	err = db.Where("live_id = ? and user_id = ? ", liveId, userInfo.UserId).First(liveRoomUser).Error
	return
}

//只有用户在直播间，才能心跳
//更新用户的心跳
//如果用户是当前直播间的主播，更新直播间心跳
func (s *Service) Heartbeat(context context.Context, liveId string, userId string) (*model.LiveEntity, error) {
	log := logger.ReqLogger(context)

	liveUser, err := s.getOrCreateLiveRoomUser(context, userId)
	if err != nil {
		log.Errorf("get live user error %v", err)
		return nil, err
	}

	if liveUser.Status != model.LiveRoomUserStatusOnline || liveUser.LiveId != liveId {
		log.Errorf("user live room (liveId: %s, status: %d), not in %s", liveUser.LiveId, liveUser.Status, liveId)
		return nil, errors.New("user not in live room")
	}

	live, err := s.getLive(context, liveId)
	if err != nil {
		log.Errorf("get live %s error %v", liveId, err)
		return live, err
	}
	go s.updateLiveUserHeartBeat(context, liveUser)
	if userId == live.AnchorId {
		go s.updateLiveHeartBeat(context, live)
	}

	return live, nil
}

func (s *Service) updateLiveHeartBeat(ctx context.Context, live *model.LiveEntity) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	updates := map[string]interface{}{
		"last_heartbeat_at": timestamp.Now(),
	}
	err := db.Model(live).Update(updates).Error
	if err != nil {
		log.Errorf("update live heart beat error: %v", err)
	}
}

func (s *Service) updateLiveUserHeartBeat(ctx context.Context, liveUser *model.LiveRoomUserEntity) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	updates := map[string]interface{}{
		"heart_beat_at": timestamp.Now(),
	}
	err := db.Model(liveUser).Update(updates).Error
	if err != nil {
		log.Errorf("update live user heart beat error: %v", err)
	}
}

func (s *Service) getLive(ctx context.Context, liveId string) (*model.LiveEntity, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	live := model.LiveEntity{}
	err := db.Where("live_id = ? ", liveId).First(&live).Error
	if err != nil {
		log.Errorf("get live error: %v", err)
		return nil, err
	}

	return &live, err
}

//查询用户在直播间内的记录
func (s *Service) getOrCreateLiveRoomUser(ctx context.Context, userId string) (*model.LiveRoomUserEntity, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	liveUser := &model.LiveRoomUserEntity{}
	result := db.Where("user_id = ?", userId).First(&liveUser)
	if result.Error != nil && !result.RecordNotFound() {
		log.Errorf("get live room user %s error %v", userId, result.Error)
		return nil, result.Error
	}

	if result.Error == nil {
		return liveUser, nil
	}

	log.Infof("create live room user %s", userId)
	liveUser = &model.LiveRoomUserEntity{
		LiveId:      "",
		UserId:      userId,
		Status:      model.LiveRoomUserStatusLeave,
		HeartBeatAt: nil,
		CreatedAt:   timestamp.Now(),
		UpdatedAt:   timestamp.Now(),
		DeletedAt:   nil,
	}
	err := db.Save(liveUser).Error
	if err != nil {
		log.Errorf("save live user error %v", err)
		return nil, err
	}

	return liveUser, nil
}

func (s *Service) SearchLive(context context.Context, keyword string, flag, pageNum, pageSize int) (lives []model.LiveEntity, totalCount int, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	lives = make([]model.LiveEntity, 0)
	if flag == 1 {
		err = db.Where(" status = ? and title like ?", model.LiveStatusOn, keyword+"%").Order("updated_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&lives).Error
		err = db.Model(&model.LiveEntity{}).Where(" status = ? and title like ?", model.LiveStatusOn, keyword+"%").Count(&totalCount).Error
	} else if flag == 2 {
		err = db.Where(" status = ? and live_id like ?", model.LiveStatusOn, keyword+"%").Order("updated_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&lives).Error
		err = db.Model(&model.LiveEntity{}).Where("status = ? and live_id like ?", model.LiveStatusOn, keyword+"%").Count(&totalCount).Error
	} else {
		err = db.Where("status = ? and anchor_id like ?", model.LiveStatusOn, keyword+"%").Order("updated_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&lives).Error
		err = db.Model(&model.LiveEntity{}).Where("status = ? and anchor_id like ?", model.LiveStatusOn, keyword+"%").Count(&totalCount).Error
	}
	if err != nil {
		log.Errorf("search live error: %v", err)
		return
	}
	return
}

func (s *Service) updateLiveStatus(context context.Context, liveId string, status int) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	live := &model.LiveEntity{}
	err = db.Where("live_id = ? ", liveId).First(live).Error
	if err != nil {
		log.Errorf("get live error: %v", err)
		return
	}
	live.Status = status
	err = db.Save(live).Error
	return
}

func (s *Service) CurrentLiveRoom(context context.Context, userId string) (liveEntity *model.LiveEntity, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	liveRoomUser := model.LiveRoomUserEntity{}
	err = db.Where(" user_id = ? and status = ?", userId, model.LiveRoomUserStatusOnline).First(&liveRoomUser).Error
	if err != nil {
		log.Errorf("find live room user error %s", err.Error())
		return nil, err
	}

	liveId := liveRoomUser.LiveId
	liveInfo, err := s.LiveInfo(context, liveId)
	if err != nil {
		log.Errorf("get live info error: %v", err)
		return nil, err
	}

	if liveInfo.Status == model.LiveStatusOn && liveInfo.AnchorId == userId {
		return liveInfo, nil
	} else {
		return nil, api.ErrNotFound
	}
}

func (s *Service) StartRelay(ctx context.Context, roomId, userId string, sid string) (err error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	live := &model.LiveEntity{}
	err = db.Where("live_id = ? and anchor_id = ?  and status = ?", roomId, userId, model.LiveStatusOn).First(live).Error
	if err != nil {
		log.Errorf("get live error: %v", err)
		return err
	}
	live.PkId = sid
	err = db.Save(live).Error
	return
}

func (s *Service) StopRelay(ctx context.Context, roomId, userId string, sid string) (err error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	live := &model.LiveEntity{}
	err = db.Where("live_id = ? and anchor_id = ?  and status = ?", roomId, userId, model.LiveStatusOn).First(live).Error
	if err != nil {
		log.Errorf("get live error: %v", err)
		return err
	}
	live.PkId = ""
	err = db.Save(live).Error
	return
}
