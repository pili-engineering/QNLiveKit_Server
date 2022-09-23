// @Author: wangsheng
// @Description:
// @File:  Client
// @Version: 1.0.0
// @Date: 2022/5/25 10:21 上午
// Copyright 2021 QINIU. All rights reserved

package pili

import (
	"fmt"
	"net/http"
	"time"

	"github.com/qiniu/go-sdk/v7/pili"

	"github.com/qbox/livekit/common/auth/qiniumac"
)

type QiniuClient struct {
	Hub string
	Ak  string
	Sk  string

	RoomTokenExpireS int64
	StreamPattern    string
	PublishDomain    string
	RtmpPlayUrl      string
	FlvPlayUrl       string
	HlsPlayUrl       string
	securityType     string //expiry, expiry_sk, none
	publishKey       string // 七牛RTC生成推流地址需要验算公共Key
	publishExpireS   int64  // 推流地址过期时间，秒
	client           *http.Client
}

const RtcHost = "https://rtc.qiniuapi.com"

func NewQiniuClient(conf Config) *QiniuClient {
	mac := &qiniumac.Mac{
		AccessKey: conf.AccessKey,
		SecretKey: []byte(conf.SecretKey),
	}
	tr := qiniumac.NewTransport(mac, nil)

	c := &QiniuClient{
		Hub:            conf.Hub,
		Ak:             conf.AccessKey,
		Sk:             conf.SecretKey,
		StreamPattern:  conf.StreamPattern,
		PublishDomain:  conf.PublishDomain,
		RtmpPlayUrl:    conf.RtmpPlayUrl,
		FlvPlayUrl:     conf.FlvPlayUrl,
		HlsPlayUrl:     conf.HlsPlayUrl,
		securityType:   conf.SecurityType,
		publishKey:     conf.PublishKey,
		publishExpireS: conf.PublishExpireS,
		client: &http.Client{
			Transport: tr,
			Timeout:   3 * time.Second,
		},
	}
	return c
}

func (c *QiniuClient) StreamPubURL(roomId string, expectAt *time.Time) (url string) {
	rtmp := pili.RTMPPublishURL(c.Hub, c.PublishDomain, c.streamName(roomId))
	var expireAt int64
	if expectAt != nil {
		expireAt = expectAt.Unix()
	} else {
		expireAt = time.Now().Add(time.Second * time.Duration(c.publishExpireS)).Unix()
	}
	url, _ = pili.SignPublishURL(rtmp, pili.SignPublishURLArgs{
		SecurityType: c.securityType,
		PublishKey:   c.publishKey,
		ExpireAt:     expireAt,
		Nonce:        0,
		AccessKey:    c.Ak,
		SecretKey:    c.Sk,
	})
	return
}

func (c *QiniuClient) StreamRtmpPlayURL(roomId string) (url string) {
	url = fmt.Sprintf(c.RtmpPlayUrl + "/" + c.Hub + "/" + c.streamName(roomId))
	return
}

func (c *QiniuClient) StreamFlvPlayURL(roomId string) (url string) {
	url = fmt.Sprintf(c.FlvPlayUrl + "/" + c.Hub + "/" + c.streamName(roomId) + ".flv")
	return
}

func (c *QiniuClient) StreamHlsPlayURL(roomId string) (url string) {
	url = fmt.Sprintf(c.HlsPlayUrl + "/" + c.Hub + "/" + c.streamName(roomId) + ".m3u8")
	return
}

func (c *QiniuClient) streamName(roomId string) string {
	return fmt.Sprintf(c.StreamPattern, roomId)
}
