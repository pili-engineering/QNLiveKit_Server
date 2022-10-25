package notify

type LikeNotifyItem struct {
	LiveId string `json:"live_id"` //直播间ID
	UserId string `json:"user_id"` // 点赞用户id
	Count  int64  `json:"count"`   //点赞数
}

type GiftNotifyItem struct {
	LiveId string `json:"live_id"` //直播间ID
	UserId string `json:"user_id"` // 发送礼物用户id
	Type   int    `json:"type"`    //礼物类型
	Amount int64  `json:"amount"`  //礼物金额
}
