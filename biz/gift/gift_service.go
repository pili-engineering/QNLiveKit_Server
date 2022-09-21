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
	Type   int    `json:"type"`
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
			Type:   req.Type,
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
			LiveId:  liveEntity.LiveId,
			Message: req.UserId + "打赏",
		}
		err = notify.SendNotifyToLive(context, anchorInfo, liveEntity, notify.ActionTypeGiftBroadcast, &notifyItem)
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
	LiveId  string
	Message string
}

func EqualsGiftRequest(req *SendGiftRequest, liveGift *model.LiveGift) bool {
	if req.UserId != liveGift.UserId || req.LiveId != liveGift.LiveID || req.Type != liveGift.Type || req.Amount != liveGift.Amount {
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
