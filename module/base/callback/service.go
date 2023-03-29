// @Author: wangsheng
// @Description:
// @File:  impl
// @Version: 1.0.0
// @Date: 2022/7/13 4:05 下午
// Copyright 2021 QINIU. All rights reserved

package callback

import (
	"context"
	"encoding/json"

	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/rpc"
)

type ICallbackService interface {
	Do(ctx context.Context, typ string, body interface{}) error
}

type CallbackService struct {
	addr   string
	client *rpc.Client
}

var callService ICallbackService = nil

func GetCallbackService() ICallbackService {
	return callService
}

type Request struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *Response) Success() bool {
	return r.Code == 0
}

func (r *Response) Error() string {
	if r.Code == 0 {
		return ""
	}

	body, _ := json.Marshal(r)
	return string(body)
}

func (s *CallbackService) Do(ctx context.Context, typ string, body interface{}) error {
	log := logger.ReqLogger(ctx)
	if len(s.addr) == 0 {
		log.Infof("no callback ")
		return nil
	}

	req := &Request{
		Type: typ,
		Body: body,
	}

	res := Response{}
	err := s.client.CallWithJSON(log, &res, s.addr, req)
	if err != nil {
		log.Errorf("callback error %s", err.Error())
		return err
	}

	if !res.Success() {
		log.Errorf("callback error %s", res.Error())
		return &res
	}

	return nil
}

const (
	TypeLiveCreated = "live_created"
	TypeLiveStarted = "live_started"
	TypeLiveStopped = "live_stopped"
	TypeLiveDeleted = "live_deleted"
	TypePKStarted   = "pk_started"
	TypePKStopped   = "pk_stopped"
	TypeGiftSend    = "gift_send" // 送出礼物会将送礼信息进行回调
)
