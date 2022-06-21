package maxim

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/rpc"
)

type Client struct {
	appId       string
	apiEndPoint string
	accessToken string
}

func NewMaximClient(imAppId string, accessToken string, endPoint string) *Client {
	return &Client{
		appId:       imAppId,
		apiEndPoint: endPoint,
		accessToken: accessToken,
	}
}

func (c *Client) defaultRpcClient() *rpc.Client {
	header := http.Header{}
	header.Set("app_id", c.appId)
	header.Set("access-token", c.accessToken)

	return rpc.NewClientHeader(header)
}

func (c *Client) rpcClientWithUserId(userId int64) *rpc.Client {
	header := http.Header{}
	header.Set("app_id", c.appId)
	header.Set("access-token", c.accessToken)
	header.Set("user_id", fmt.Sprintf("%d", userId))

	return rpc.NewClientHeader(header)
}

// RegisterUser /user/register/v2 注册用户,返回注册的用户id
func (c *Client) RegisterUser(ctx context.Context, username, password string) (int64, error) {
	log := logger.ReqLogger(ctx)
	url := c.apiEndPoint + "/user/register/v2"
	var user map[string]string = map[string]string{
		"username": username,
		"password": password,
	}

	resp := RegisterUserResponse{}
	err := c.defaultRpcClient().CallWithJSON(log, &resp, url, user)
	if err != nil {
		log.Errorf("register user error %s", err.Error())
		return 0, api.ErrInternal
	}

	if resp.IsSuccess() {
		return resp.Data.UserId, nil
	}
	log.Errorf("register user error %s", resp.Error())
	if resp.Code != 10004 {
		return 0, api.ErrInternal
	}

	//用户已存在的错误
	userId, err := c.GetUserId(ctx, username)
	if err != nil {
		log.Errorf("get user id for %s error %s", username, err)
		return 0, api.ErrInternal
	}

	_, err = c.UpdateUserPassword(ctx, userId, password)
	if err != nil {
		log.Errorf("update password for %d error %s", userId, err)
		return 0, api.ErrInternal
	}

	return userId, nil
}

func (c *Client) CreateChatroom(ctx context.Context, owner int64, name string) (int64, error) {
	log := logger.ReqLogger(ctx)
	url := c.apiEndPoint + "/group/create"
	var req = map[string]interface{}{
		"name": name,
		"type": 2,
	}

	resp := CreateChatRoomResponse{}
	err := c.rpcClientWithUserId(owner).CallWithJSON(log, &resp, url, req)
	if err != nil {
		log.Errorf("create group error %s", err.Error())
		return 0, api.ErrInternal
	}

	if resp.IsSuccess() {
		return resp.Data.GroupId, nil
	} else {
		log.Errorf("create group error %s", resp.Error())
		return 0, api.ErrInternal
	}
}

func (c *Client) GetUserId(ctx context.Context, username string) (int64, error) {
	log := logger.ReqLogger(ctx)
	url := c.apiEndPoint + "/roster/name?username=" + username

	resp := GetUserResponse{}
	err := c.defaultRpcClient().GetCall(log, &resp, url)
	if err != nil {
		log.Errorf("get user id error %s", err.Error())
		return 0, api.ErrInternal
	}

	if resp.IsSuccess() {
		return resp.Data.UserId, nil
	} else {
		log.Errorf("get user id for %s error %s", username, resp.Error())
		return 0, api.ErrInternal
	}
}

func (c *Client) UpdateUserPassword(ctx context.Context, userId int64, password string) (bool, error) {
	log := logger.ReqLogger(ctx)
	url := c.apiEndPoint + "/user/change_password_admin"
	var req = map[string]interface{}{
		"password": password,
	}

	resp := UpdatePasswordResponse{}
	err := c.rpcClientWithUserId(userId).CallWithJSON(log, &resp, url, req)
	if err != nil {
		log.Errorf("update password error %s", err.Error())
		return false, api.ErrInternal
	}

	if resp.IsSuccess() {
		return resp.Data, nil
	} else {
		log.Errorf("update password for %d error %s", userId, resp.Error())
		return false, api.ErrInternal
	}
}
