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
	"github.com/qbox/livekit/utils/rpc"
	"github.com/qbox/livekit/utils/uuid"
)

type SendGiftRequest struct {
	LiveId string `json:"live_id"`
	GiftId int    `json:"gift_id"`
	Amount int    `json:"amount"`
}

func (s *Service) SendGift(context context.Context, req *SendGiftRequest, userId string) (*SendGiftResponse, error) {
	log := logger.ReqLogger(context)
	bizId := uuid.Gen()
	liveEntity, err := live.GetService().LiveInfo(context, req.LiveId)
	if err != nil {
		log.Errorf("find live error %s", err.Error())
		return nil, api.ErrorGiftPay
	}
	anchorInfo, err := user.GetService().FindUser(context, liveEntity.AnchorId)
	if err != nil {
		log.Errorf("get anchor info for %s error %s", liveEntity.AnchorId, err.Error())
		return nil, api.ErrorGiftPay
	}
	liveGift := &model.LiveGift{
		LiveId:   req.LiveId,
		UserId:   userId,
		BizId:    bizId,
		GiftId:   req.GiftId,
		Amount:   req.Amount,
		AnchorId: anchorInfo.UserId,
	}
	err = SaveLiveGift(context, liveGift)
	if err != nil {
		log.Errorf("save live gift error %s", err.Error())
		return nil, api.ErrorGiftPay
	}

	payReq := &PayGiftRequest{
		LiveId:   req.LiveId,
		UserId:   userId,
		BizId:    bizId,
		GiftId:   req.GiftId,
		Amount:   req.Amount,
		AnchorId: anchorInfo.UserId,
	}
	payResp := PayGiftResponse{}
	url := s.GiftHost
	err = rpc.DefaultClient.CallWithJSON(log, &payResp, url, payReq)
	if err != nil {
		return nil, err
	}

	sResp := &SendGiftResponse{
		LiveId:   req.LiveId,
		UserId:   userId,
		BizId:    bizId,
		GiftId:   req.GiftId,
		Amount:   req.Amount,
		AnchorId: anchorInfo.UserId,
		Status:   payResp.Status,
	}
	err = s.UpdateGiftStatus(context, bizId, payResp.Status)
	if err != nil {
		log.Errorf("update gift status error %s", err.Error())
		//该错误没有返回
	}
	if payResp.Status == model.SendGiftLiveStatusFailure {
		return sResp, api.ErrorGiftPayFromBiz
	}
	notifyItem := BroadcastGiftNotifyItem{
		LiveId: liveEntity.LiveId,
		UserId: userId,
		GiftId: req.GiftId,
		Amount: req.Amount,
	}
	err = notify.SendNotifyToLive(context, anchorInfo, liveEntity, notify.ActionTypeGiftNotify, &notifyItem)
	if err != nil {
		log.Errorf("send notify to live %s error %s", liveEntity.LiveId, err.Error())
		return sResp, api.ErrorGiftPayIMMessage
	}
	return sResp, nil
}

type PayGiftRequest struct {
	BizId    string `json:"biz_id"`
	UserId   string `json:"user_id"`
	LiveId   string `json:"live_id"`
	AnchorId string `json:"anchor_id"`
	GiftId   int    `json:"gift_id"`
	Amount   int    `json:"amount"`
}

type PayGiftResponse struct {
	api.Response
	Status int
}

type SendGiftResponse struct {
	BizId    string `json:"biz_id"`
	UserId   string `json:"user_id"`
	LiveId   string `json:"live_id"`
	AnchorId string `json:"anchor_id"`
	GiftId   int    `json:"gift_id"`
	Amount   int    `json:"amount"`
	Status   int    `json:"status"`
}

type BroadcastGiftNotifyItem struct {
	UserId string `json:"user_id"`
	LiveId string `json:"live_id"`
	GiftId int    `json:"gift_id"`
	Amount int    `json:"amount"`
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

func (s *Service) UpdateGiftStatus(context context.Context, bizId string, status int) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err := db.Model(&model.LiveGift{}).Where("biz_id = ?", bizId).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SearchGiftByAnchorId(context context.Context, anchorId string, pageNum, pageSize int) (liveGifts []*model.LiveGift, totalCount int, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	liveGifts = make([]*model.LiveGift, 0)
	err = db.Model(&model.LiveGift{}).Where("anchor_id = ?", anchorId).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&liveGifts).Error
	err = db.Model(&model.LiveGift{}).Where("anchor_id = ?", anchorId).Count(&totalCount).Error
	if err != nil {
		log.Errorf("SearchGiftByAnchorId %v", err)
		return nil, 0, err
	}
	return
}

func (s *Service) SearchGiftByLiveId(context context.Context, liveId string, pageNum, pageSize int) (liveGifts []*model.LiveGift, totalCount int, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	liveGifts = make([]*model.LiveGift, 0)
	err = db.Model(&model.LiveGift{}).Where("live_id = ?", liveId).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&liveGifts).Error
	err = db.Model(&model.LiveGift{}).Where("live_id = ?", liveId).Count(&totalCount).Error
	if err != nil {
		log.Errorf("SearchGiftByLiveId %v", err)
		return nil, 0, err
	}
	return
}

func (s *Service) SearchGiftByUserId(context context.Context, userId string, pageNum, pageSize int) (liveGifts []*model.LiveGift, totalCount int, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	liveGifts = make([]*model.LiveGift, 0)
	err = db.Model(&model.LiveGift{}).Where("user_id = ?", userId).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&liveGifts).Error
	err = db.Model(&model.LiveGift{}).Where("user_id = ?", userId).Count(&totalCount).Error
	if err != nil {
		log.Errorf("SearchGiftByUserId %v", err)
		return nil, 0, err
	}
	return
}
