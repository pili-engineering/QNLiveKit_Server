// @Author: wangsheng
// @Description:
// @File:  live_timeout
// @Version: 1.0.0
// @Date: 2022/6/1 2:28 下午
// Copyright 2021 QINIU. All rights reserved

package impl

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/timestamp"

	"github.com/qbox/livekit/utils/logger"
)

const LiveUserHeartBeatTime = 5 * time.Second
const LiveUserHeartBeatTimeout = 4 * LiveUserHeartBeatTime

const LiveRoomHeartBeatTime = 5 * time.Second
const LiveRoomHeartBeatTimeout = 10 * time.Minute

func (s *service.Service) TimeoutLiveUser(ctx context.Context, now time.Time) {
	log := logger.ReqLogger(ctx)
	defer func() {
		if e := recover(); e != nil {
			const size = 16 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Error("TimeoutLiveUser  panic: ", e, fmt.Sprintf("\n%s", buf))
		}
	}()

	liveUsers, _ := s.listAllTimeoutLiveUsers(ctx)
	if len(liveUsers) == 0 {
		return
	}

	for _, liveUser := range liveUsers {
		s.LeaveLiveRoom(ctx, liveUser.LiveId, liveUser.UserId)
	}
}

func (s *service.Service) listAllTimeoutLiveUsers(ctx context.Context) ([]*model.LiveRoomUserEntity, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	liveUsers := make([]*model.LiveRoomUserEntity, 0)
	result := db.Find(&liveUsers, "status = ? and heart_beat_at < ?", model.LiveRoomUserStatusOnline, timestamp.Now().Add(-1*LiveUserHeartBeatTimeout))
	if result.Error != nil {
		log.Errorf("list timeout live users error %v", result.Error)
		return nil, result.Error
	}

	return liveUsers, nil
}

func (s *service.Service) TimeoutLiveRoom(ctx context.Context, now time.Time) {
	log := logger.ReqLogger(ctx)
	defer func() {
		if e := recover(); e != nil {
			const size = 16 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Error("TimeoutLiveRoom  panic: ", e, fmt.Sprintf("\n%s", buf))
		}
	}()

	liveRooms, _ := s.listAllTimeoutLiveRooms(ctx)
	if len(liveRooms) == 0 {
		return
	}

	for _, liveRoom := range liveRooms {
		s.StopLive(ctx, liveRoom.LiveId, liveRoom.AnchorId)
	}
}

func (s *service.Service) listAllTimeoutLiveRooms(ctx context.Context) ([]*model.LiveEntity, error) {
	log := logger.ReqLogger(ctx)
	lives := make([]*model.LiveEntity, 0)
	db := mysql.GetLive(log.ReqID())
	result := db.Find(&lives, "status = ? and last_heartbeat_at < ?", model.LiveStatusOn, timestamp.Now().Add(-LiveRoomHeartBeatTimeout))
	if result.Error != nil {
		log.Errorf("list timeout live room error %v", result.Error)
		return nil, result.Error
	}

	return lives, nil
}
