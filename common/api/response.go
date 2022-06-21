// @Author: wangsheng
// @Description:
// @File:  response
// @Version: 1.0.0
// @Date: 2022/5/19 3:55 下午
// Copyright 2021 QINIU. All rights reserved

package api

import "encoding/json"

type Response struct {
	RequestId string `json:"request_id"` //请求ID
	Code      int    `json:"code"`       //错误码，0 成功，其他失败
	Message   string `json:"message"`    //错误信息
}

func (r *Response) Error() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func SuccessResponse(reqId string) Response {
	return Response{
		RequestId: reqId,
		Code:      0,
		Message:   "success",
	}
}
