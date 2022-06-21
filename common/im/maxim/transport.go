// @Author: wangsheng
// @Description:
// @File:  transport
// @Version: 1.0.0
// @Date: 2022/5/20 4:48 下午
// Copyright 2021 QINIU. All rights reserved

package maxim

import (
	"net/http"
)

type Transport struct {
	appId       string
	accessToken string
	Transport   http.RoundTripper
}

func NewTransport(appId string, accessToken string, transport http.RoundTripper) *Transport {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return &Transport{
		appId:       appId,
		accessToken: accessToken,
		Transport:   transport,
	}
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("app_id", t.appId)
	req.Header.Set("access-token", t.accessToken)
	return t.Transport.RoundTrip(req)
}

func (t *Transport) NestedObject() interface{} {
	return t.Transport
}
