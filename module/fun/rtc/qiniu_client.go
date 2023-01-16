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
	"github.com/qiniu/go-sdk/v7/rtc"

	"github.com/qbox/livekit/core/module/account"
	"github.com/qbox/livekit/utils/qiniumac"
)

type QiniuClient struct {
	AppId  string
	Ak     string
	Sk     string
	client *http.Client
}

func NewQiniuClient(conf Config) *QiniuClient {
	c := &QiniuClient{
		AppId: conf.AppId,
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

func (c *QiniuClient) setupAccount() error {
	// 未配置，不需要设置账号
	if service == nil {
		return nil
	}

	if account.AccessKey() == "" {
		return fmt.Errorf("no account info")
	}

	service.Ak = account.AccessKey()
	service.Sk = account.SecretKey()

	mac := &qiniumac.Mac{
		AccessKey: account.AccessKey(),
		SecretKey: []byte(account.SecretKey()),
	}
	tr := qiniumac.NewTransport(mac, nil)
	service.client = &http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}
	return nil
}

func (c *QiniuClient) RtcAppId() string {
	return c.AppId
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
