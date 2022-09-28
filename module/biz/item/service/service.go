package service

import (
	"context"

	"github.com/qbox/livekit/biz/model"
)

var Instance IItemService

type IItemService interface {
	AddItems(ctx context.Context, liveId string, items []*model.ItemEntity) error
	DelItems(ctx context.Context, liveId string, items []string) error
	ListItems(ctx context.Context, liveId string, showOffline bool) ([]*model.ItemEntity, error)
	UpdateItemInfo(ctx context.Context, liveId string, item *model.ItemEntity) error
	UpdateItemExtends(ctx context.Context, liveId string, itemId string, extends model.Extends) error

	UpdateItemStatus(ctx context.Context, liveId string, statuses []*model.ItemStatus) error
	UpdateItemOrder(ctx context.Context, liveId string, orders []*model.ItemOrder) error
	UpdateItemOrderSingle(ctx context.Context, liveId string, itemId string, from, to uint) error

	SetDemonstrateItem(ctx context.Context, liveId string, itemId string) error
	DelDemonstrateItem(ctx context.Context, liveId string) error
	GetDemonstrateItem(ctx context.Context, liveId string) (*model.ItemEntity, error)
	GetLiveItem(ctx context.Context, liveId string, itemId string) (*model.ItemEntity, error)
	IsDemonstrateItem(ctx context.Context, liveId, itemId string) (bool, error)

	StartRecordVideo(ctx context.Context, liveId string, itemId string) error
	StopRecordVideo(ctx context.Context, liveId string, demonId int) (*model.ItemDemonstrateRecord, error)
	GetRecordVideo(ctx context.Context, demonId uint) (*model.ItemDemonstrateRecord, error)
	UpdateRecordVideo(ctx context.Context, itemLog *model.ItemDemonstrateRecord) error
	GetListRecordVideo(ctx context.Context, liveId string, itemId string) (*model.ItemDemonstrateRecord, error)
	GetListLiveRecordVideo(ctx context.Context, liveId string) ([]*model.ItemDemonstrateRecord, error)
	GetPreviousItem(ctx context.Context, liveId string) (*int, error)
	DelRecordVideo(ctx context.Context, liveId string, demonItem []uint) error
	SaveRecordVideo(ctx context.Context, liveId, itemId string) error
	UpdateItemRecord(ctx context.Context, demonId uint, liveId string, itemId string) error
	DeleteItemRecord(ctx context.Context, demonId uint, liveId string, itemId string) error
}
