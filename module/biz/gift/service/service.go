package service

import (
	"context"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/module/biz/gift/internal/impl"
)

type Service interface {
	SaveGiftEntity(context context.Context, entity *model.GiftEntity) error
	DeleteGiftEntity(context context.Context, giftId int) error
	GetListGiftEntity(context context.Context, typeId int) ([]*model.GiftEntity, error)
	SendGift(context context.Context, req *impl.SendGiftRequest, userId string) (*impl.SendGiftResponse, error)
	SearchGiftByLiveId(context context.Context, liveId string, pageNum, pageSize int) (liveGifts []*model.LiveGift, totalCount int, err error)
	SearchGiftByAnchorId(context context.Context, anchorId string, pageNum, pageSize int) (liveGifts []*model.LiveGift, totalCount int, err error)
	SearchGiftByUserId(context context.Context, userId string, pageNum, pageSize int) (liveGifts []*model.LiveGift, totalCount int, err error)
	UpdateGiftStatus(context context.Context, bizId string, status int) error
}
