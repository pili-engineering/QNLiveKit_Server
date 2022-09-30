package impl

import (
	"context"
	"fmt"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/module/biz/mic/service"
	"github.com/qbox/livekit/module/fun/rtc"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/logger"
)

type ServiceImpl struct {
}

var instance = &ServiceImpl{}

func GetInstance() *ServiceImpl {
	return instance
}

func (s *ServiceImpl) KickUser(context context.Context, userId, liveId string) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	userMic1 := &model.LiveMicEntity{}
	result := db.Model(userMic1).Where("live_id = ? and user_id = ? and type = ?", liveId, userId, "mic").First(userMic1)
	if result.Error != nil {
		if !result.RecordNotFound() {
			log.Errorf("kick user error: %v", result.Error)
			return result.Error
		}
	}
	if userMic1.Id != 0 {
		userMic1.Status = model.LiveRoomUserMicForbidden
	}
	userMic2 := &model.LiveMicEntity{}
	result = db.Model(userMic2).Where("live_id = ? and user_id = ? and type = ?", liveId, userId, "camera").First(userMic2)
	if result.Error != nil {
		if !result.RecordNotFound() {
			log.Errorf("kick user error: %v", result.Error)
			return result.Error
		}
	}
	if userMic2.Id != 0 {
		userMic2.Status = model.LiveRoomUserMicForbidden
	}
	return
}

func (s *ServiceImpl) UpMic(context context.Context, req *service.Request, userId string) (rtcToken string, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	userMic := &model.LiveMicEntity{}
	result := db.Model(userMic).Where("live_id = ? and user_id = ?", req.LiveId, userId).First(userMic)
	if result.Error != nil {
		if result.RecordNotFound() {
			userMic.LiveId = req.LiveId
			userMic.UserId = userId
			userMic.Mic = req.Mic
			userMic.Camera = req.Camera
			userMic.Status = model.LiveRoomUserMicStatusJoin
			userMic.Extends = req.Extends
			result = db.Create(userMic)
			if result.Error != nil {
				log.Errorf("up mic error: %v", result.Error)
				err = result.Error
				return
			}
		} else {
			log.Errorf("up mic error: %v", result.Error)
			err = result.Error
			return
		}
	} else {
		if userMic.Status == model.LiveRoomUserMicForbidden {
			err = fmt.Errorf("user mic forbidden")
			return
		}
		userMic.Status = model.LiveRoomUserMicStatusJoin
		err = db.Save(userMic).Error
	}

	rtcService := rtc.GetService()
	rtcToken = rtcService.GetRoomToken(userId, req.LiveId)
	return
}

func (s *ServiceImpl) DownMic(context context.Context, req *service.Request, userId string) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	userMic := &model.LiveMicEntity{}
	result := db.Model(userMic).Where("live_id = ? and user_id = ?", req.LiveId, userId).First(userMic)
	if result.Error != nil {
		if result.RecordNotFound() {
			log.Errorf("down mic error: %v", result.Error)
			return
		} else {
			log.Errorf("down mic error: %v", result.Error)
			return result.Error
		}
	} else {
		userMic.Status = model.LiveRoomUserMicStatusLeave
		err = db.Save(userMic).Error
	}
	return
}

func (s *ServiceImpl) DownMicManual(context context.Context, liveId, userId string) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	userMic := &model.LiveMicEntity{}
	result := db.Model(userMic).Where("live_id = ? and user_id = ? ", liveId, userId).First(userMic)
	if result.Error != nil {
		if result.RecordNotFound() {
			log.Errorf("down mic error: %v", result.Error)
			return
		} else {
			log.Errorf("down mic error: %v", result.Error)
			return result.Error
		}
	} else {
		userMic.Status = model.LiveRoomUserMicStatusLeave
		err = db.Save(userMic).Error
	}
	return
}

func (s *ServiceImpl) ForbidMic(context context.Context, liveId, userId string) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	userMic := &model.LiveMicEntity{}
	result := db.Model(userMic).Where("live_id = ? and user_id = ?", liveId, userId).First(userMic)
	if result.Error != nil {
		if result.RecordNotFound() {
			return
		} else {
			log.Errorf("forbid mic error: %v", result.Error)
			return result.Error
		}
	} else {
		userMic.Status = model.LiveRoomUserMicForbidden
		err = db.Save(userMic).Error
	}
	return
}

func (s *ServiceImpl) UnForbidMic(context context.Context, liveId, userId string) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	userMic := &model.LiveMicEntity{}
	result := db.Model(userMic).Where("live_id = ? and user_id = ?", liveId, userId).First(userMic)
	if result.Error != nil {
		if result.RecordNotFound() {
			return
		} else {
			log.Errorf("release mic error: %v", result.Error)
			return result.Error
		}
	} else {
		userMic.Status = model.LiveRoomUserMicStatusLeave
		err = db.Save(userMic).Error
	}
	return
}

func (s *ServiceImpl) UserStatus(context context.Context, liveId, userId string) (status int, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	userMic := &model.LiveMicEntity{}
	result := db.Model(userMic).Where("live_id = ? and user_id = ?", liveId, userId).First(userMic)
	if result.Error != nil {
		if result.RecordNotFound() {
			status = -1
			return
		} else {
			log.Errorf("get user status error: %v", result.Error)
			status = -1
			err = result.Error
			return
		}
	}
	status = userMic.Status
	return
}

func (s *ServiceImpl) LiveMicList(context context.Context, liveId string) (mics []model.LiveMicEntity, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	mics = make([]model.LiveMicEntity, 0)
	err = db.Where("live_id = ? and status = ?", liveId, model.LiveRoomUserMicStatusJoin).Find(&mics).Error
	return
}

func (s *ServiceImpl) UpdateMicExtends(context context.Context, liveId, userId string, extends model.Extends) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	userMic := &model.LiveMicEntity{}
	result := db.Model(userMic).Where("live_id = ? and user_id = ?", liveId, userId).First(userMic)
	if result.Error != nil {
		if result.RecordNotFound() {
			return
		} else {
			log.Errorf("update user mic error: %v", result.Error)
			return result.Error
		}
	}
	userMic.Extends = extends
	err = db.Save(userMic).Error
	return
}

func (s *ServiceImpl) SwitchUserMic(context context.Context, liveId, userId, tp string, flag bool) (err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	userMic := &model.LiveMicEntity{}
	result := db.Model(userMic).Where("live_id = ? and user_id = ?", liveId, userId).First(userMic)
	if result.Error != nil {
		if result.RecordNotFound() {
			return
		} else {
			log.Errorf("switch user mic error: %v", result.Error)
			return result.Error
		}
	}
	if tp == "mic" {
		userMic.Mic = flag
	} else if tp == "camera" {
		userMic.Camera = flag
	} else {
		return
	}
	err = db.Save(userMic).Error
	return
}
