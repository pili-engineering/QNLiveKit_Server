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
//

//

//
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
