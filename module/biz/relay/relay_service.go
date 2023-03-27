// @Author: wangsheng
// @Description:
// @File:  impl
// @Version: 1.0.0
// @Date: 2022/5/26 10:13 上午
// Copyright 2021 QINIU. All rights reserved

package relay

import (
	"context"
	"github.com/qbox/livekit/biz/notify"
	"github.com/qbox/livekit/module/base/user"

	"github.com/qbox/livekit/core/module/uuid"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/live"
	"github.com/qbox/livekit/module/store/mysql"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

var relayService IRelayService = &RelayService{}

func GetRelayService() IRelayService {
	return relayService
}

type StartRelayParams struct {
	InitUserId string //发起方用户ID
	InitRoomId string //
	RecvRoomId string
	RecvUserId string
	Extends    model.Extends
}

type IRelayService interface {
	//开始跨房
	StartRelay(ctx context.Context, params *StartRelayParams) (*model.RelaySession, error)

	//上报自己跨房完成
	//appId
	//userId 上报用户ID
	//sid 跨房会话ID
	ReportRelayStarted(ctx context.Context, userId string, sid string) (*model.RelaySession, error)

	//获取跨房会话
	//appId 直播应用ID
	//sid 跨房会话ID
	GetRelaySession(ctx context.Context, sid string) (*model.RelaySession, error)

	//停止跨房
	//appId
	//userId 上报用户ID
	//sid 跨房会话ID
	StopRelay(ctx context.Context, userId string, sid string) error

	//获取跨房的目的房间
	//appId
	//userId 上报用户ID
	//sid 跨房会话ID
	GetRelayRoom(ctx context.Context, userId string, sid string) (*model.RelaySession, string, error)

	//更新扩展信息
	UpdateRelayExtends(ctx context.Context, sid string, extends model.Extends) error
}

type RelayService struct {
}

func (s *RelayService) StartRelay(ctx context.Context, params *StartRelayParams) (*model.RelaySession, error) {
	log := logger.ReqLogger(ctx)

	if err := s.checkCanStartRelay(ctx, params.InitRoomId, params.InitUserId); err != nil {
		return nil, err
	}

	if err := s.checkCanStartRelay(ctx, params.RecvRoomId, params.RecvUserId); err != nil {
		return nil, err
	}

	relaySession, err := s.createRelaySession(ctx, params)
	if err != nil {
		log.Errorf("create relay session error %v", err)
		return nil, err
	}

	liveService := live.GetService()
	err = liveService.StartRelay(ctx, params.InitRoomId, params.InitUserId, relaySession.SID)
	if err != nil {
		return nil, err
	}

	err = liveService.StartRelay(ctx, params.RecvRoomId, params.RecvUserId, relaySession.SID)
	if err != nil {
		return nil, err
	}

	return relaySession, nil
}

func (s *RelayService) checkCanStartRelay(ctx context.Context, roomId string, userId string) error {
	log := logger.ReqLogger(ctx)
	liveService := live.GetService()
	liveEntity, err := liveService.LiveInfo(ctx, roomId)
	if err != nil {
		log.Errorf("get live room %s error %v", roomId, err)
		return rest.ErrNotFound
	}
	if liveEntity.Status != model.LiveStatusOn || liveEntity.AnchorId != userId {
		log.Errorf("live room anchorId (%s, %s) not equal to userId (%s)", liveEntity.LiveId, liveEntity.AnchorId, userId)
		return rest.ErrNotFound
	}

	if liveEntity.PkId != "" {
		log.Errorf("live room (%s) is relaying", roomId)
		return rest.ErrBadRequest.WithMessage("live room already in relaying")
	}

	return nil
}

func (s *RelayService) createRelaySession(ctx context.Context, params *StartRelayParams) (*model.RelaySession, error) {
	log := logger.ReqLogger(ctx)

	relaySession := model.RelaySession{
		SID:        uuid.Gen(),
		InitUserId: params.InitUserId,
		InitRoomId: params.InitRoomId,
		RecvUserId: params.RecvUserId,
		RecvRoomId: params.RecvRoomId,
		Extends:    params.Extends,
		Status:     model.RelaySessionStatusAgreed,
		StartAt:    nil,
		StopAt:     nil,
		CreatedAt:  timestamp.Timestamp{},
		UpdatedAt:  timestamp.Timestamp{},
	}

	db := mysql.GetLive(log.ReqID())
	if err := db.Create(&relaySession).Error; err != nil {
		log.Errorf("create relay session error %v", err)
		return nil, rest.ErrInternal
	}

	return &relaySession, nil
}

func (s *RelayService) ReportRelayStarted(ctx context.Context, userId string, sid string) (*model.RelaySession, error) {
	log := logger.ReqLogger(ctx)

	relaySession, _, err := s.GetRelayRoom(ctx, userId, sid)
	if err != nil {
		return nil, err
	}

	if relaySession.InitUserId == userId {
		if relaySession.Status == model.RelaySessionStatusAgreed {
			relaySession.Status = model.RelaySessionStatusInitSuccess
		} else {
			relaySession.Status = model.RelaySessionStatusSuccess
		}
	} else {
		if relaySession.Status == model.RelaySessionStatusAgreed {
			relaySession.Status = model.RelaySessionStatusRecvSuccess
		} else {
			relaySession.Status = model.RelaySessionStatusSuccess
		}
	}

	err = s.updateRelaySessionStatus(ctx, relaySession)
	if err != nil {
		log.Errorf("update relay session error %v", err)
		return nil, err
	} else {
		return relaySession, nil
	}
}

func (s *RelayService) GetRelaySession(ctx context.Context, sid string) (*model.RelaySession, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLiveReadOnly(log.ReqID())

	relaySession := model.RelaySession{}
	result := db.First(&relaySession, "sid = ?", sid)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, rest.ErrNotFound
		} else {
			log.Errorf("find relay session error %v", result.Error)
			return nil, rest.ErrInternal
		}
	}

	return &relaySession, nil
}

func (s *RelayService) StopRelay(ctx context.Context, userId string, sid string) error {
	log := logger.ReqLogger(ctx)

	relaySession, _, err := s.GetRelayRoom(ctx, userId, sid)
	if err != nil {
		log.Errorf("get relay room error %v", err)
		return err
	}

	var relayUser, relayRoom, room string
	if relaySession.InitUserId == userId {
		room = relaySession.InitRoomId

		relayRoom = relaySession.RecvRoomId
		relayUser = relaySession.RecvUserId
	} else {
		room = relaySession.RecvRoomId

		relayRoom = relaySession.InitRoomId
		relayUser = relaySession.InitUserId
	}
	relaySession.Status = model.RelaySessionStatusStopped
	s.updateRelaySessionStatus(ctx, relaySession)

	liveService := live.GetService()

	liveService.StopRelay(ctx, room, userId, sid)
	liveService.StopRelay(ctx, relayRoom, relayUser, sid)

	return nil
}

// 更新跨房PK 会话的状态
func (s *RelayService) updateRelaySessionStatus(ctx context.Context, relaySession *model.RelaySession) error {
	log := logger.ReqLogger(ctx)

	updates := map[string]interface{}{
		"status":     relaySession.Status,
		"updated_at": timestamp.Now(),
	}

	db := mysql.GetLive(log.ReqID())
	result := db.Model(relaySession).Update(updates)
	if result.Error != nil {
		log.Errorf("update relay session status error %v", result.Error)
		return result.Error
	}
	return nil
}

// 获取用户当前的跨房房间，前提
// 1，用户当前正在直播
// 2，用户当前正在跨房
func (s *RelayService) GetRelayRoom(ctx context.Context, userId string, sid string) (*model.RelaySession, string, error) {
	log := logger.ReqLogger(ctx)

	relaySession, err := s.GetRelaySession(ctx, sid)
	if err != nil {
		log.Errorf("get relay session (%s) error %v", sid, err)
		return nil, "", err
	}

	if relaySession.IsStopped() {
		log.Errorf("relaySession %s stopped", sid)
		return nil, "", rest.ErrNotFound
	}

	liveService := live.GetService()
	liveRoom, err := liveService.CurrentLiveRoom(ctx, userId)
	if err != nil {
		log.Errorf("get current live room error %v", err)
		return nil, "", err
	}

	relayRoom := ""
	if userId == relaySession.InitUserId && liveRoom.LiveId == relaySession.InitRoomId {
		relayRoom = relaySession.RecvRoomId
	} else if userId == relaySession.RecvUserId && liveRoom.LiveId == relaySession.RecvRoomId {
		relayRoom = relaySession.InitRoomId
	} else {
		log.Errorf("user room (%s) not in relay session (%s)", liveRoom.LiveId, relaySession.SID)
		return nil, "", rest.ErrNotFound
	}

	return relaySession, relayRoom, nil
}

func (s *RelayService) UpdateRelayExtends(ctx context.Context, sid string, extends model.Extends) error {
	log := logger.ReqLogger(ctx)

	relaySession, err := s.GetRelaySession(ctx, sid)
	if err != nil {
		log.Errorf("get relay session (%s) error %v", sid, err)
		return err
	}

	relaySession.Extends = model.CombineExtends(relaySession.Extends, extends)
	return s.updateRelayExtends(ctx, relaySession)
}

func (s *RelayService) updateRelayExtends(ctx context.Context, relaySession *model.RelaySession) error {
	log := logger.ReqLogger(ctx)

	updates := map[string]interface{}{
		"extends":    relaySession.Extends,
		"updated_at": timestamp.Now(),
	}
	db := mysql.GetLive(log.ReqID())
	// 查询出更新之前的值
	oldRelaySession := &model.RelaySession{}
	db.Model(&model.RelaySession{}).Where("id = ?", relaySession.ID).Find(oldRelaySession)
	result := db.Model(relaySession).Update(updates)
	if result.Error != nil {
		log.Errorf("update relay session extends error %v", result.Error)
		return result.Error
	}
	// 拓展字段更新完后的hook
	go afterUpdateRelayExtendsHandler(ctx, relaySession, oldRelaySession)
	return nil
}

// afterUpdateRelayExtendsHandler
//
//	@Description: 拓展字段更新完后的操作，将更新的字段发给双方房间
//	@args ctx
//	@args relaySession
func afterUpdateRelayExtendsHandler(ctx context.Context, relaySession *model.RelaySession, oldRelaySession *model.RelaySession) {
	log := logger.ReqLogger(ctx)
	liveUserEntityList, err := user.GetService().ListUser(ctx, []string{relaySession.InitUserId, relaySession.RecvUserId})
	if err != nil {
		log.Errorf("cannot queried the relaySession【%v】，errInfo：【%v】", relaySession, err.Error())
		return
	}
	// 通过用户的的id
	// 要发送的房间，pk的两方房间
	userId2LiveEntityMap, err := user.GetService().FindLiveByPkIdList(ctx, relaySession.SID)
	if err != nil {
		log.Errorf("cannot queried the LiveRoom【%v】，errInfo：【%v】", relaySession.InitRoomId, err.Error())
		return
	}
	data := &extendNotifyData{
		SID:           relaySession.SID,
		InitRoomId:    relaySession.InitRoomId,
		RecvRoomId:    relaySession.RecvRoomId,
		UpdateExtends: getUpdateFieldMap(relaySession, oldRelaySession),
	}
	for _, userEntity := range liveUserEntityList {
		liveEntity, ok := userId2LiveEntityMap[userEntity.UserId]
		if ok {
			err = notify.SendNotifyToLive(ctx, userEntity, liveEntity, notify.ActionTypeExtendsNotify, data)
			if err != nil {
				log.Errorf("cannot send notify, LiveUserEntity【%v】|data【%v】|errInfo【%v】", userEntity, data, err.Error())
			} else {
				log.Infof("【pk_extends_notify】send notify success, LiveUserEntity【%v】| data:【%v】", userEntity, data)
			}
		} else {
			log.Errorf("cannot find liveEntity, userEntity【%v】", userEntity)
		}
	}
}

// getUpdateFieldMap 返回更新了的字段
func getUpdateFieldMap(newSession *model.RelaySession, oldSession *model.RelaySession) map[string]string {
	if oldSession == nil {
		return newSession.Extends
	}
	updateFieldMap := make(map[string]string)
	// 返回空map。表示更新空值
	if newSession == nil {
		return updateFieldMap
	}
	for k, v := range newSession.Extends {
		if oldSession.Extends[k] != v {
			updateFieldMap[k] = v
		}
	}
	return updateFieldMap
}

type extendNotifyData struct {
	SID           string            `json:"sid"`          // PK 会话ID
	InitRoomId    string            `json:"init_room_id"` // 发送发直播间ID
	RecvRoomId    string            `json:"recv_room_id"` // 接收方直播间ID
	UpdateExtends map[string]string `json:"extends"`      // 更新了的字段
}
