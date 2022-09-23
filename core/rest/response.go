package rest

type Response struct {
	RequestId string      `json:"request_id"`     //请求ID
	Code      int         `json:"code"`           //错误码，0 成功，其他失败
	Message   string      `json:"message"`        //错误信息
	Data      interface{} `json:"data,omitempty"` //传递的数据
}
