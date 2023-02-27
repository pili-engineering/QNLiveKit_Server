package service

import (
	"context"

	"github.com/qbox/livekit/biz/model"
)

type IUserService interface {
	FindUser(ctx context.Context, userId string) (*model.LiveUserEntity, error)
	FindOrCreateUser(ctx context.Context, userId string) (*model.LiveUserEntity, error)
	ListUser(ctx context.Context, userIds []string) ([]*model.LiveUserEntity, error)
	ListImUser(ctx context.Context, imUserIds []int64) ([]*model.LiveUserEntity, error)
	CreateUser(ctx context.Context, user *model.LiveUserEntity) error
	UpdateUserInfo(ctx context.Context, user *model.LiveUserEntity) error
	// FindLiveByPkIdList 根据PK会话查询直播间信息
	FindLiveByPkIdList(ctx context.Context, pkIdList ...string) (liveRoomUser *[]model.LiveEntity, err error)
}
