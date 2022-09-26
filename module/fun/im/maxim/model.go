// @Author: wangsheng
// @Description:
// @File:  model
// @Version: 1.0.0
// @Date: 2022/5/20 5:05 下午
// Copyright 2021 QINIU. All rights reserved

package maxim

import "encoding/json"

type Response struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (r *Response) IsSuccess() bool {
	return r.Code == 200
}

func (r *Response) Error() string {
	data, _ := json.Marshal(r)
	return string(data)
}

type RegisterUserResponse struct {
	Response
	Data struct {
		UserId int64 `json:"user_id"`
	} `json:"data"`
}

type CreateChatRoomResponse struct {
	Response
	Data struct {
		GroupId int64 `json:"group_id"`
	} `json:"data"`
}

type GetUserResponse struct {
	Response
	Data struct {
		UserId int64 `json:"user_id"`
	} `json:"data"`
}

type UpdatePasswordResponse struct {
	Response
	Data bool `json:"data"`
}

type CommonResponse struct {
	Response
	Data bool `json:"data"`
}
