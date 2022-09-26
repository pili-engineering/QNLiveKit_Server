package live

import (
	"context"
	"errors"
	"time"

	"github.com/qbox/livekit/biz/admin"
	"github.com/qbox/livekit/core/module/uuid"
	"github.com/qbox/livekit/module/base/callback"
	"github.com/qbox/livekit/module/fun/im"
	"github.com/qbox/livekit/module/fun/pili"
	"github.com/qbox/livekit/module/fun/rtc"
	"github.com/qbox/livekit/module/store/mysql"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/user"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

type IService interface {
	CreateLive(context context.Context, req *CreateLiveRequest) (live *model.LiveEntity, err error)

	GetLiveAuthor(ctx context.Context, liveId string) (*model.LiveUserEntity, error)

	DeleteLive(context context.Context, liveId string, anchorId string) (err error)

	StartLive(context context.Context, liveId string, anchorId string) (roomToken string, err error)

	StopLive(context context.Context, liveId string, anchorId string) (err error)

	AdminStopLive(ctx context.Context, liveId string, reason string, adminId string) error

	LiveInfo(context context.Context, liveId string) (live *model.LiveEntity, err error)

	LiveListAnchor(context context.Context, pageNum, pageSize int, anchorId string) (lives []model.LiveEntity, totalCount int, err error)

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

	CheckLiveAnchor(ctx context.Context, liveId string, userId string) error

	UpdateLiveRelatedReview(context context.Context, liveId string, latest *int) (err error)

	AddLike(ctx context.Context, liveId string, userId string, count int64) (my, total int64, err error)

	FlushCacheLikes(ctx context.Context)
}

type Service struct {
}

var service IService = &Service{}

func GetService() IService {
	return service
}

type CreateLiveRequest struct {
	AnchorId        string               `json:"anchor_id"`
	Title           string               `json:"title"`
	Notice          string               `json:"notice"`
	CoverUrl        string               `json:"cover_url"`
	StartAt         timestamp.Timestamp  `json:"start_at"`
	EndAt           timestamp.Timestamp  `json:"end_at"`
	PublishExpireAt *timestamp.Timestamp `json:"publish_expire_at"`
	Extends         model.Extends        `json:"extends" gorm:"type:varchar(512)"`
}

func (s *Service) CreateLive(context context.Context, req *CreateLiveRequest) (live *model.LiveEntity, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	liveId := uuid.Gen()

	liveUser, err := user.GetService().FindUser(context, req.AnchorId)
	if err != nil {
		log.Errorf("create live failed, user not found, userId: %s, err: %v", req.AnchorId, err)
		return
	}
	imClient := im.GetService()
	chatroom, err := imClient.CreateChatroom(context, liveUser.ImUserid, liveId)
	if err != nil {
		log.Errorf("create chatroom failed, err: %v", err)
		return
	}
	exp := req.PublishExpireAt
	var url string
	if exp != nil && exp.After(time.Now()) {
		url = pili.GetService().StreamPubURL(liveId, &exp.Time)
	} else {
		url = pili.GetService().StreamPubURL(liveId, nil)
	}

	live = &model.LiveEntity{
		LiveId:      liveId,
		Title:       req.Title,
		Notice:      req.Notice,
		CoverUrl:    req.CoverUrl,
		Extends:     req.Extends,
		AnchorId:    req.AnchorId,
		Status:      model.LiveStatusPrepare,
		PkId:        "",
		OnlineCount: 0,
		StartAt:     &req.StartAt, //timestamp.Now(),
		EndAt:       req.EndAt,    //timestamp.Now(),
		ChatId:      chatroom,
		PushUrl:     url,
		RtmpPlayUrl: pili.GetService().StreamRtmpPlayURL(liveId),
		FlvPlayUrl:  pili.GetService().StreamFlvPlayURL(liveId),
		HlsPlayUrl:  pili.GetService().StreamHlsPlayURL(liveId),
	}
	err = db.Create(live).Error
	if err == nil {
		err = admin.GetCensorService().CreateCensorJob(context, live)
		if err != nil {
			log.Errorf("create censor job  failed, err: %v", err)
		}
		go callback.GetCallbackService().Do(context, callback.TypeLiveCreated, live)
	}
	return
}

func (s *Service) GetLiveAuthor(ctx context.Context, liveId string) (*model.LiveUserEntity, error) {
	log := logger.ReqLogger(ctx)
	liveInfo, err := s.LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("get live %s error %v", liveId, err)
		return nil, err
	}

	userInfo, err := user.GetService().FindUser(ctx, liveInfo.AnchorId)
	if err != nil {
		log.Errorf("get user %s error %v", liveInfo.AnchorId, err)
		return nil, err
	}

	return userInfo, nil
}

func (s *Service) DeleteLive(context context.Context, liveId string, anchorId string) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err = db.Delete(&model.LiveEntity{}, "live_id = ? and anchor_id = ? ", liveId, anchorId).Error

	if err == nil {
		body := map[string]string{
			"live_id": liveId,
		}
		go callback.GetCallbackService().Do(context, callback.TypeLiveDeleted, body)
	}
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
	if live.Status == model.LiveStatusOff {
		err = errors.New("live status error")
		return
	}
	if live.AnchorId != anchorId {
		err = errors.New("user not anchor")
		return
	}

	//判断主播不在其他直播间
	liveUser, err := s.getOrCreateLiveRoomUser(context, anchorId)
	if err != nil {
		return "", err
	}

	live.Status = model.LiveStatusOn
	now := timestamp.Now()
	live.StartAt = &now
	live.LastHeartbeatAt = timestamp.Now()
	live.UpdatedAt = timestamp.Now()
	err = db.Save(live).Error
	roomToken = rtcClient.GetRoomToken(anchorId, liveId)

	liveUser.Status = model.LiveRoomUserStatusOnline
	liveUser.LiveId = liveId
	liveUser.UpdatedAt = now
	liveUser.HeartBeatAt = &now

	db.Save(liveUser)

	body := map[string]string{
		"live_id": liveId,
	}
	go callback.GetCallbackService().Do(context, callback.TypeLiveStarted, body)

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

	if err == nil {
		body := map[string]string{
			"live_id": liveId,
		}
		go callback.GetCallbackService().Do(context, callback.TypeLiveStopped, body)
	}

	err = admin.GetCensorService().StopCensorJob(context, liveId)
	if err != nil {
		log.Errorf("stop censor job failed, err: %v", err)
		return err
	}
	return
}

func (s *Service) AdminStopLive(ctx context.Context, liveId string, reason string, adminId string) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	live, err := s.LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("LiveInfo error:%v", err)
		return err
	}
	if live.Status != model.LiveStatusOn {
		err = errors.New("live status error")
		return err
	}

	now := timestamp.Now()

	live.Status = model.LiveStatusOff
	live.EndAt = now

	live.StopReason = reason
	live.StopUserId = adminId
	live.StopAt = &now

	err = db.Save(live).Error

	if err == nil {
		body := map[string]string{
			"live_id": liveId,
		}
		go callback.GetCallbackService().Do(ctx, callback.TypeLiveStopped, body)
	}

	return nil
}

func (s *Service) LiveInfo(context context.Context, liveId string) (live *model.LiveEntity, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	live = &model.LiveEntity{}
	err = db.Where("live_id = ? ", liveId).First(live).Error
	return
}

func (s *Service) LiveListAnchor(context context.Context, pageNum, pageSize int, anchorId string) (lives []model.LiveEntity, totalCount int, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	lives = make([]model.LiveEntity, 0)
	err = db.Where("anchor_id = ?", anchorId).Order("start_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&lives).Error
	err = db.Model(&model.LiveEntity{}).Where("anchor_id = ?", anchorId).Count(&totalCount).Error
	return
}

func (s *Service) LiveList(context context.Context, pageNum, pageSize int) (lives []model.LiveEntity, totalCount int, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	lives = make([]model.LiveEntity, 0)
	err = db.Where("status = ? or status = ?", model.LiveStatusOn, model.LiveStatusPrepare).Order("status desc").Order("updated_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&lives).Error
	err = db.Model(&model.LiveEntity{}).Where("status =  ? or status = ?", model.LiveStatusOn, model.LiveStatusPrepare).Count(&totalCount).Error
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

func (s *Service) UpdateLiveRelatedReview(context context.Context, liveId string, latest *int) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	unreviewCount, err := admin.GetCensorService().GetUnreviewCount(context, liveId)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{}
	updates["unreview_censor_count"] = unreviewCount
	if latest != nil {
		updates["last_censor_time"] = *latest
	}
	result := db.Model(&model.LiveEntity{}).Where("live_id = ? ", liveId).Update(updates)
	if result.Error != nil {
		log.Errorf("update live about censor information error %v", result.Error)
		return api.ErrDatabase
	} else {
		return nil
	}
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

func (s *Service) LeaveLiveRoom(context context.Context, liveId string, userId string) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	liveRoomUser := &model.LiveRoomUserEntity{}
	result := db.Where("live_id = ? and user_id = ? ", liveId, userId).First(liveRoomUser)
	if result.Error != nil {
		log.Errorf("find live rooom user error %s", result.Error.Error())
		return api.ErrDatabase
	} else {
		liveRoomUser.LiveId = ""
		liveRoomUser.Status = model.LiveRoomUserStatusLeave
		liveRoomUser.UpdatedAt = timestamp.Now()
		if err := db.Save(liveRoomUser).Error; err != nil {
			log.Errorf("update live room user error %s", err.Error())
			return api.ErrDatabase
		}
	}

	// 如果是主播离开房间了，取消直播间当前讲解商品
	liveEntity, err := s.getLive(context, liveId)
	if err != nil {
		log.Errorf("get live error %s", err.Error())
		return api.ErrDatabase
	}

	if liveEntity != nil && liveEntity.AnchorId == userId {
		itemService := GetItemService()
		err = itemService.DelDemonstrateItem(context, liveId)
		if err != nil {
			log.Errorf("delete demonstrate for live %s error %s", liveId, err.Error())
		}
	}
	return nil
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

func (s *Service) CheckLiveAnchor(ctx context.Context, liveId string, userId string) error {
	log := logger.ReqLogger(ctx)

	liveEntity, err := s.LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("get live info error %v", err)
		return err
	}

	if liveEntity.AnchorId != userId {
		log.Errorf("anchor not match liveAnchor(%s), userId(%s)", liveEntity.AnchorId, userId)
		return api.ErrNotFound
	}

	return nil
}

func (s *Service) AddLike(ctx context.Context, liveId string, userId string, count int64) (my, total int64, err error) {
	log := logger.ReqLogger(ctx)
	liveInfo, err := s.LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("get liveInfo info failed, err: %v", err)
		return 0, 0, api.ErrDatabase
	}

	if liveInfo.AnchorId == userId {
		log.Errorf("anchor can not like self")
		return 0, 0, api.ErrInvalidArgument
	}

	liveUser, err := s.getOrCreateLiveRoomUser(ctx, userId)
	if err != nil {
		log.Errorf("get live user error %v", err)
		return 0, 0, err
	}

	if liveUser.Status != model.LiveRoomUserStatusOnline || liveUser.LiveId != liveId {
		log.Errorf("user live room (liveId: %s, status: %d), not in %s", liveUser.LiveId, liveUser.Status, liveId)
		return 0, 0, errors.New("user not in live room")
	}

	my, total, err = s.cacheLike(ctx, liveId, userId, count)
	if err != nil {
		return 0, 0, err
	}

	return my, total, nil
}
