// @Author: wangsheng
// @Description:
// @File:  Client
// @Version: 1.0.0
// @Date: 2022/5/25 10:21 上午
// Copyright 2021 QINIU. All rights reserved

package rtc

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/pili"
	"github.com/qiniu/go-sdk/v7/rtc"

	"github.com/qbox/livekit/common/auth/qiniumac"
)

type QiniuClient struct {
	RoomTokenExpireS int64
	RtcPlayBackUrl   string
	Hub              string
	StreamPattern    string
	PublishUrl       string
	PublishDomain    string
	RtmpPlayUrl      string
	FlvPlayUrl       string
	HlsPlayUrl       string
	AppId            string
	Ak               string
	Sk               string
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
		AppId:            conf.AppId,
		Ak:               conf.AccessKey,
		Sk:               conf.SecretKey,
		RoomTokenExpireS: conf.RoomTokenExpireS,
		RtcPlayBackUrl:   conf.RtcPlayBackUrl,
		Hub:              conf.Hub,
		StreamPattern:    conf.StreamPattern,
		PublishUrl:       conf.PublishUrl,
		PublishDomain:    conf.PublishDomain,
		RtmpPlayUrl:      conf.RtmpPlayUrl,
		FlvPlayUrl:       conf.FlvPlayUrl,
		HlsPlayUrl:       conf.HlsPlayUrl,
		client: &http.Client{
			Transport: tr,
			Timeout:   3 * time.Second,
		},
		securityType:   conf.SecurityType,
		publishKey:     conf.PublishKey,
		publishExpireS: conf.PublishExpireS,
	}
	return c
}

type RoomAccess struct {
	AppID      string `json:"appId"`
	RoomName   string `json:"roomName"`
	UserID     string `json:"userId"`
	ExpireAt   int64  `json:"expireAt"`
	Permission string `json:"permission"`
	//Privileges string `json:"privileges,omitempty"`
	//Scenario   int    `json:"scenario,omitempty"`
}

func (c *QiniuClient) GetRoomToken(userId, roomName string) string {
	roomAccess := &RoomAccess{
		AppID:      c.AppId,
		RoomName:   roomName,
		UserID:     userId,
		ExpireAt:   time.Now().Unix() + 24*3600,
		Permission: "user",
	}

	return c.signToken(roomAccess)
}

func (c *QiniuClient) GetRelayToken(userId, roomName string) string {
	relayAccess := &RoomAccess{
		AppID:      c.AppId,
		RoomName:   roomName,
		UserID:     userId,
		ExpireAt:   time.Now().Unix() + 24*3600,
		Permission: "user",
	}
	return c.signToken(relayAccess)
}

func (c *QiniuClient) signToken(roomAccess interface{}) string {
	data, _ := json.Marshal(roomAccess)
	buf := make([]byte, base64.URLEncoding.EncodedLen(len(data)))
	base64.URLEncoding.Encode(buf, data)

	hmacsha1 := hmac.New(sha1.New, []byte(c.Sk))
	hmacsha1.Write(buf)
	sign := hmacsha1.Sum(nil)

	encodedSign := base64.URLEncoding.EncodeToString(sign)
	token := c.Ak + ":" + encodedSign + ":" + string(buf)
	return token
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

func (c *QiniuClient) ListUser(roomId string) (res []string, err error) {
	signer := &auth.Credentials{
		AccessKey: c.Ak,
		SecretKey: []byte(c.Sk),
	}
	manager := rtc.NewManager(signer)
	users, err := manager.ListUser(c.AppId, roomId)
	if err != nil {
		return nil, err
	} else {
		res = make([]string, 0, len(users))
		for _, u := range users {
			res = append(res, u.UserID)
		}
		return
	}
}

func (c *QiniuClient) Online(userId, roomId string) bool {
	result := make(chan bool)
	go func() {
		users, err := c.ListUser(roomId)
		if err != nil {
			result <- false
		}
		for _, id := range users {
			if id == userId {
				result <- true
			}
		}
		result <- false
	}()
	select {
	case res := <-result:
		return res
	case <-time.After(5 * time.Second):
		return false
	}
}
