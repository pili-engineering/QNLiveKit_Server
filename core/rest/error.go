package rest

import (
	"net/http"
)

type Error struct {
	StatusCode int    `json:"-"`          //http 状态码
	RequestId  string `json:"request_id"` //请求ID
	Code       int    `json:"code"`       //错误码，0 成功，其他失败
	Message    string `json:"message"`    //错误信息
}

func IsNotFoundError(err error) bool {
	if e1, ok := err.(*Error); ok {
		return e1.StatusCode == ErrNotFound.StatusCode && e1.Code == ErrNotFound.Code
	} else {
		return false
	}
}

func (e *Error) Error() string {
	return ""
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) WithMessage(message string) *Error {
	e1 := *e
	e1.Message = message
	return &e1
}

func (e *Error) WithStatusCode(statusCode int) *Error {
	e1 := *e
	e1.StatusCode = statusCode
	return &e1
}

func (e *Error) WithRequestId(requestId string) *Error {
	e1 := *e
	e1.RequestId = requestId
	return &e1
}

func (e *Error) WithMessageAndCode(code int, message string, requestId string) *Error {
	e1 := *e
	e1.Code = code
	e1.RequestId = requestId
	if message != "" {
		e1.Message = message
	}
	return &e1
}

// 基础错误
var ErrNotFound = &Error{StatusCode: http.StatusNotFound, Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}
var ErrBadRequest = &Error{StatusCode: http.StatusBadRequest, Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)}
var ErrUnauthorized = &Error{StatusCode: http.StatusUnauthorized, Code: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
var ErrForbidden = &Error{StatusCode: http.StatusForbidden, Code: http.StatusForbidden, Message: http.StatusText(http.StatusForbidden)}
var ErrTimeout = &Error{StatusCode: http.StatusRequestTimeout, Code: http.StatusRequestTimeout, Message: http.StatusText(http.StatusRequestTimeout)}
var ErrInternal = &Error{StatusCode: http.StatusInternalServerError, Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}

var ErrTokenExpired = &Error{StatusCode: http.StatusUnauthorized, Code: 499, Message: "Your token is expired"}
var ErrAlreadyExist = &Error{StatusCode: http.StatusBadRequest, Code: 10001, Message: "Already existed"}
var ErrGiftPay = &Error{StatusCode: http.StatusOK, Code: 20002, Message: "PayGift Failure"}
