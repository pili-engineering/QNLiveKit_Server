package gift

import (
	"context"
	"github.com/qbox/livekit/biz/live"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/notify"
	"github.com/qbox/livekit/biz/user"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/utils/logger"
)

type SendGiftRequest struct {
	BizId  string `json:"biz_id"`
	UserId string `json:"user_id"`
	LiveId string `json:"live_id"`
	GiftId int    `json:"gift_id"`
	Amount int    `json:"amount"`
	Redo   bool   `json:"redo"`
}

func SendGift(context context.Context, req SendGiftRequest) error {
	log := logger.ReqLogger(context)
	gift, err := GetGiftByBizId(context, req.BizId)
	if err != nil {
		return err
	}
	if gift == nil {
		err = SaveLiveGift(context, &model.LiveGift{
			LiveID: req.LiveId,
			UserId: req.UserId,
			BizId:  req.BizId,
			GiftId: req.GiftId,
			Amount: req.Amount,
		})
		if err != nil {
			return api.ErrDatabase
		}

		liveEntity, err := live.GetService().LiveInfo(context, req.LiveId)
		if err != nil {
			log.Errorf("find live error %s", err.Error())
			return api.ErrDatabase
		}
		anchorInfo, err := user.GetService().FindUser(context, liveEntity.AnchorId)
		if err != nil {
			log.Errorf("get anchor info for %s error %s", liveEntity.AnchorId, err.Error())
		}

		notifyItem := BroadcastGiftNotifyItem{
			LiveId: liveEntity.LiveId,
			UserId: req.UserId,
			GiftId: req.GiftId,
			Amount: req.Amount,
		}
		err = notify.SendNotifyToLive(context, anchorInfo, liveEntity, notify.ActionTypeGiftNotify, &notifyItem)
		if err != nil {
			log.Errorf("send notify to live %s error %s", liveEntity.LiveId, err.Error())
		}
	} else {
		if req.Redo == false {
			return api.ErrorGiftBizIdRepeatedWrong
		} else {
			equals := EqualsGiftRequest(&req, gift)
			if equals {
				return nil
			}
			return api.ErrorGiftRedoDataInconsistency
		}
	}
	return nil
}

type BroadcastGiftNotifyItem struct {
	UserId string `json:"user_id"`
	LiveId string `json:"live_id"`
	GiftId int    `json:"gift_id"`
	Amount int    `json:"amount"`
}

func EqualsGiftRequest(req *SendGiftRequest, liveGift *model.LiveGift) bool {
	if req.UserId != liveGift.UserId || req.LiveId != liveGift.LiveID || req.GiftId != liveGift.GiftId || req.Amount != liveGift.Amount {
		return false
	}
	return true

}

func GetGiftByBizId(context context.Context, bizId string) (liveGift *model.LiveGift, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	result := db.Model(&model.LiveGift{}).Find(liveGift, "biz_id = ?", bizId)
	if result.Error != nil {
		if !result.RecordNotFound() {
			log.Errorf("get gift by biz_id  error: %v", result.Error)
			return nil, result.Error
		}
	}
	return
}

func SaveLiveGift(context context.Context, liveGift *model.LiveGift) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLiveReadOnly(log.ReqID())
	err := db.Model(&model.LiveGift{}).Save(liveGift).Error
	if err != nil {
		return api.ErrDatabase
	}
	return nil
}

func UpdateGiftStatus(context context.Context, giftId int, status int) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err := db.Model(&model.LiveGift{}).Where("gift_id = ?", giftId).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func GetGiftByAnchorId(context context.Context, anchorId int) ([]*model.LiveGift, error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	gifts := make([]*model.LiveGift, 0)
	err := db.Model(&model.LiveGift{}).Find(&gifts, "anchor_id = ?", anchorId).Error
	if err != nil {
		return nil, err
	}
	return gifts, nil
}

func GetGiftByLiveId(context context.Context, liveId int) ([]*model.LiveGift, error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	gifts := make([]*model.LiveGift, 0)
	err := db.Model(&model.LiveGift{}).Find(&gifts, "live_id = ?", liveId).Error
	if err != nil {
		return nil, err
	}
	return gifts, nil
}

func GetGiftByUserId(context context.Context, userId int) ([]*model.LiveGift, error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	gifts := make([]*model.LiveGift, 0)
	err := db.Model(&model.LiveGift{}).Find(&gifts, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return gifts, nil
}
