// @Author: wangsheng
// @Description:
// @File:  errors
// @Version: 1.0.0
// @Date: 2022/5/19 3:55 下午
// Copyright 2021 QINIU. All rights reserved

package api

import "net/http"

func ErrorWithRequestId(reqID string, err error) *Response {
	resp := &Response{
		RequestId: reqID,
	}

	if r, ok := err.(*Response); ok {
		resp.Code = r.Code
		resp.Message = r.Message
	} else {
		resp.Code = 501
		resp.Message = err.Error()
	}

	return resp
}

func Error(reqID string, code int, message string) *Response {
	return &Response{
		RequestId: reqID,
		Code:      code,
		Message:   message,
	}
}

func IsNotFoundError(err error) bool {
	switch e := err.(type) {
	case *Response:
		return e.Code == ErrorCodeNotFound

	default:
		return false
	}
}

const (
	ErrorCodeNotFound        = http.StatusNotFound
	ErrorCodeInvalidArgument = http.StatusBadRequest
	ErrorCodeBadToken        = http.StatusUnauthorized
	ErrorCodeInternal        = http.StatusInternalServerError

	ErrorCodeAlreadyExisted = 600
	ErrorCodeBadStatus      = 601
	ErrorCodeDatabase       = 602
	ErrorCodeTokenExpired   = 499

	ErrorCodeUserAlreadyExisted = 10001 //用户已经存在

	ErrorCodeLiveItemExceed = 20001 //直播间商品数量超过限制

	ErrorCodeLoginWrong = 30001 //admin登录密码或用户名错误

	ErrorCodeGiftPay           = 200002
	ErrorCodeGiftPayRequestBiz = 200003
)

var ErrInvalidArgument = &Response{Code: ErrorCodeInvalidArgument, Message: "The arguments you provide is invalid."}
var ErrNotFound = &Response{Code: ErrorCodeNotFound, Message: "Not found"}
var ErrAlreadyExist = &Response{Code: ErrorCodeAlreadyExisted, Message: "Already existed"}
var ErrBadToken = &Response{Code: ErrorCodeBadToken, Message: "Your authorization token is invalid"}
var ErrTokenExpired = &Response{Code: ErrorCodeTokenExpired, Message: "Your token is expired"}
var ErrInternal = &Response{Code: ErrorCodeInternal, Message: "Internal error"}
var ErrDatabase = &Response{Code: ErrorCodeDatabase, Message: "Database error"}
var ErrStatus = &Response{Code: ErrorCodeBadStatus, Message: "cant operate on this status"}

var ErrCodeLiveItemExceed = &Response{Code: ErrorCodeLiveItemExceed, Message: "items exceed in live room"}

var ErrorLoginWrong = &Response{Code: ErrorCodeLoginWrong, Message: "username or password wrong，Login Failure"}

var ErrorGiftPay = &Response{Code: ErrorCodeGiftPay, Message: "PayGift Failure: Internal error"}

var ErrorGiftPayRequestBiz = &Response{Code: ErrorCodeGiftPayRequestBiz, Message: "PayGift Failure: Request Biz Server Error"}

func FailGiftResponse(code int, msg string) *Response {
	if msg == "" {
		msg = "PayGift Failure"
	}
	return &Response{
		Code:    code,
		Message: msg,
	}
}
