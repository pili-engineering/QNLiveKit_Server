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

	job := SkipIfStillRunning()(cron.FuncJob(cmd))

	return s.cron.AddJob(spec, job)
}

// isSingleTaskNode 判断自己是否需要单节点运行
func (s *Service) isSingleTaskNode() bool {
	return node.NodeId() == s.SingleTaskNode
}

type JobWrapper func(cron.Job) cron.Job

// SkipIfStillRunning skips an invocation of the Job if a previous invocation is
// still running. It logs skips to the given logger at Info level.
func SkipIfStillRunning() JobWrapper {
	return func(j cron.Job) cron.Job {
		var ch = make(chan struct{}, 1)
		ch <- struct{}{}
		return cron.FuncJob(func() {
			select {
			case v := <-ch:
				defer func() { ch <- v }()
				j.Run()
			default:
				// skip
			}
		})
	}
}
