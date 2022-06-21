// @Author: wangsheng
// @Description:
// @File:  service
// @Version: 1.0.0
// @Date: 2022/5/26 10:13 上午
// Copyright 2021 QINIU. All rights reserved

package live

import (
	"context"

	"github.com/qbox/livekit/common/mysql"

	"github.com/qbox/livekit/utils/uuid"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
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

	liveService := GetService()
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
	liveService := GetService()
	liveEntity, err := liveService.LiveInfo(ctx, roomId)
	if err != nil {
		log.Errorf("get live room %s error %v", roomId)
		return api.ErrNotFound
	}
	if liveEntity.Status != model.LiveStatusOn || liveEntity.AnchorId != userId {
		log.Errorf("live room anchorId (%s, %s) not equal to userId (%s)", liveEntity.LiveId, liveEntity.AnchorId, userId)
		return api.ErrNotFound
	}

	if liveEntity.RelaySessionId != "" {
		log.Errorf("live room (%s) is relaying", roomId)
		return api.Error(log.ReqID(), api.ErrorCodeBadStatus, "live room already in relaying")
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
		return nil, api.ErrDatabase
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
			return nil, api.ErrNotFound
		} else {
			log.Errorf("find relay session error %v", result.Error)
			return nil, api.ErrDatabase
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

	liveService := GetService()

	liveService.StopRelay(ctx, room, userId, sid)
	liveService.StopRelay(ctx, relayRoom, relayUser, sid)

	return nil
}

//更新跨房PK 会话的状态
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

//获取用户当前的跨房房间，前提
//1，用户当前正在直播
//2，用户当前正在跨房
func (s *RelayService) GetRelayRoom(ctx context.Context, userId string, sid string) (*model.RelaySession, string, error) {
	log := logger.ReqLogger(ctx)

	relaySession, err := s.GetRelaySession(ctx, sid)
	if err != nil {
		log.Errorf("get relay session (%s) error %v", sid, err)
		return nil, "", err
	}

	if relaySession.IsStopped() {
		log.Errorf("relaySession %s stopped", sid)
		return nil, "", api.ErrNotFound
	}

	liveService := GetService()
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
		log.Errorf("user room (%s) not in relay session (%s)", liveRoom.LiveId, relaySession.ID)
		return nil, "", api.ErrNotFound
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
	result := db.Model(relaySession).Update(updates)
	if result.Error != nil {
		log.Errorf("update relay session extends error %v", result.Error)
		return result.Error
	}
	return nil
}
