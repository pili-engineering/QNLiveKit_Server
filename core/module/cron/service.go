// @Author: wangsheng
// @Description:
// @File:  cron
// @Version: 1.0.0
// @Date: 2022/6/1 2:32 下午
// Copyright 2021 QINIU. All rights reserved

package cron

import (
	"gopkg.in/robfig/cron.v2"

	"github.com/qbox/livekit/core/module/node"
)

var instance *Service

type Service struct {
	cron           *cron.Cron
	SingleTaskNode int64
}

func newService(singleTaskNode int64) *Service {
	return &Service{
		cron:           cron.New(),
		SingleTaskNode: singleTaskNode,
	}
}

func (s *Service) StartCron() {
	s.cron.Start()
}

func (s *Service) StopCron() error {
	s.cron.Stop()
	return nil
}

// AddFunc 增加所有节点都会运行的任务
func (s *Service) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, cmd)
}

// AddSingleTaskFunc 增加单节点单线程运行的任务
func (s *Service) AddSingleTaskFunc(spec string, cmd func()) (cron.EntryID, error) {
	if !s.isSingleTaskNode() {
		return 0, nil
	}

	return s.cron.AddFunc(spec, cmd)
}

// isSingleTaskNode 判断自己是否需要单节点运行
func (s *Service) isSingleTaskNode() bool {
	return node.NodeId() == s.SingleTaskNode
}

// Run 运行 cronjob
//func Run() {
//	c := cron.New()
//
//	liveService := live.GetService()
//
//	// 定时老化直播间，单节点执行
//	c.AddFunc("0/3 * * * * ?", func() {
//		if !isSingleTaskNode() {
//			return
//		}
//
//		now := time.Now()
//		nowStr := now.Format(timestamp.TimestampFormatLayout)
//		log := logger.New("TimeoutLiveRoom")
//		log.WithFields(map[string]interface{}{"start": nowStr})
//
//		ctx := context.Background()
//		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)
//
//		liveService.TimeoutLiveRoom(ctx, now)
//	})
//
//	// 定时老化直播间用户，单节点执行
//	c.AddFunc("0/3 * * * * ?", func() {
//		if !isSingleTaskNode() {
//			return
//		}
//
//		now := time.Now()
//		nowStr := now.Format(timestamp.TimestampFormatLayout)
//
//		log := logger.New("TimeoutLiveUser")
//		log.WithFields(map[string]interface{}{"start": nowStr})
//
//		ctx := context.Background()
//		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)
//
//		liveService.TimeoutLiveUser(ctx, now)
//	})
//
//	// 上报直播间信息，单节点执行
//	c.AddFunc("0 0 2 * * ?", func() {
//		if !isSingleTaskNode() {
//			return
//		}
//
//		now := time.Now()
//		nowStr := now.Format(timestamp.TimestampFormatLayout)
//
//		log := logger.New("ReportOnlineMessage")
//		log.WithFields(map[string]interface{}{"start": nowStr})
//
//		ctx := context.Background()
//		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)
//
//		report.GetService().ReportOnlineMessage(ctx)
//	})
//
//	// 上报本节点的API 监控信息，所有节点都要执行
//	c.AddFunc("0/5 * * * * ?", func() {
//		now := time.Now()
//		nowStr := now.Format(timestamp.TimestampFormatLayout)
//
//		log := logger.New("ReportApiMonitor")
//		log.WithFields(map[string]interface{}{"start": nowStr})
//
//		ctx := context.Background()
//		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)
//
//		monitor.ReportMonitorItems(ctx)
//	})
//
//	// 每秒统计缓存中的直播间点赞，写入DB
//	// 因为存在补数据，这里不用每秒任务
//	if isSingleTaskNode() {
//		log := logger.New("FlushCacheLikes")
//
//		ctx := context.Background()
//		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)
//
//		go liveService.FlushCacheLikes(ctx)
//	}
//	c.Start()
//}
