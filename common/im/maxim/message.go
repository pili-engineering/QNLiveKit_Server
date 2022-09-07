package maxim

import (
	"context"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/rpc"
)

type ContentType int32

const (
	ContentTypeText     ContentType = 0
	ContentTypeImage    ContentType = 1
	ContentTypeAudio    ContentType = 2
	ContentTypeVideo    ContentType = 3
	ContentTypeFile     ContentType = 4
	ContentTypeLocation ContentType = 5
	ContentTypeCommand  ContentType = 6
	ContentTypeForward  ContentType = 7
)

func (t ContentType) IsValid() bool {
	return int32(t) >= 0 && int32(t) <= 7
}

type TargetType int32

const (
	TargetTypeUser  TargetType = 1
	TargetTypeGroup TargetType = 2
)

func (t TargetType) IsValid() bool {
	return t == TargetTypeUser || t == TargetTypeGroup
}

func (c *Client) SendCommandMessageToGroup(ctx context.Context, fromUserId int64, toGroupId int64, content string) error {
	req := SendMessageRequest{
		TransactionId: 0,
		Type:          TargetTypeGroup,
		FromUserId:    fromUserId,
		Targets:       []int64{toGroupId},
		ContentType:   ContentTypeCommand,
		Content:       content,
	}

	return c.sendMessage(ctx, fromUserId, &req)
}

func (c *Client) SendCommandMessageToUser(ctx context.Context, fromUserId int64, toUserId int64, content string) error {
	req := SendMessageRequest{
		TransactionId: 0,
		Type:          TargetTypeUser,
		FromUserId:    fromUserId,
		Targets:       []int64{toUserId},
		ContentType:   ContentTypeCommand,
		Content:       content,
	}

	return c.sendMessage(ctx, fromUserId, &req)
}

type SendMessageRequest struct {
	TransactionId int64       `json:"transaction_id,omitempty"`
	Type          TargetType  `json:"type"`
	FromUserId    int64       `json:"from_user_id,omitempty"`
	Targets       []int64     `json:"targets"`
	ContentType   ContentType `json:"content_type"`
	Content       string      `json:"content"`
	Attachment    string      `json:"attachment"`
	Ext           string      `json:"ext,omitempty"`
}

func (c *Client) sendMessage(ctx context.Context, fromUserId int64, req *SendMessageRequest) error {
	log := logger.ReqLogger(ctx)
	url := c.apiEndPoint + "/message/send"

	var rpcClient *rpc.Client
	if fromUserId == 0 {
		rpcClient = c.defaultRpcClient()
	} else {
		rpcClient = c.rpcClientWithUserId(fromUserId)
	}

	resp := CommonResponse{}

	err := rpcClient.CallWithJSON(log, &resp, url, req)
	if err != nil {
		log.Errorf("send message error %s", err.Error())
		return api.ErrInternal
	}

	if resp.IsSuccess() {
		return nil
	} else {
		log.Errorf("send message failed %s", resp.Error())
		return api.ErrInternal
	}
}
