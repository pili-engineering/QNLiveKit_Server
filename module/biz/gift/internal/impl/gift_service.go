package impl

import (
	"context"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/notify"
	"github.com/qbox/livekit/core/module/uuid"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/live"
	"github.com/qbox/livekit/module/base/user"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/rpc"
)

type SendGiftRequest struct {
	LiveId string `json:"live_id"`
	GiftId int    `json:"gift_id"`
	Amount int    `json:"amount"`
}

func (s *ServiceImpl) SendGift(context context.Context, req *SendGiftRequest, userId string) (*SendGiftResponse, error) {
	log := logger.ReqLogger(context)
	bizId := uuid.Gen()
	liveEntity, err := live.GetService().LiveInfo(context, req.LiveId)
	if err != nil {
		log.Errorf("find live error %s", err.Error())
		return nil, rest.ErrGiftPay
	}
	liveGift := &model.LiveGift{
		LiveId:   req.LiveId,
		UserId:   userId,
		BizId:    bizId,
		GiftId:   req.GiftId,
		Amount:   req.Amount,
		AnchorId: liveEntity.AnchorId,
	}
	err = SaveLiveGift(context, liveGift)
	if err != nil {
		log.Errorf("save live gift error %s", err.Error())
		return nil, rest.ErrGiftPay
	}

	payReq := &PayGiftRequest{
		LiveId:   req.LiveId,
		UserId:   userId,
		BizId:    bizId,
		GiftId:   req.GiftId,
		Amount:   req.Amount,
		AnchorId: liveEntity.AnchorId,
	}
	payResp := PayGiftResponse{}
	url := s.GiftAddr
	err = rpc.DefaultClient.CallWithJSON(log, &payResp, url, payReq)
	if err != nil {
		log.Errorf("send gift error %s", err.Error())
		err2 := s.UpdateGiftStatus(context, bizId, model.SendGiftStatusFailure)
		if err2 != nil {
			log.Errorf("update gift status error %s", err2.Error())
		}
		return nil, rest.ErrGiftPay
	}
	if payResp.Code != 0 {
		log.Errorf("send gift request return code not 0 ")
		err = s.UpdateGiftStatus(context, bizId, model.SendGiftStatusFailure)
		if err != nil {
			log.Errorf("update gift status error %s", err.Error())
		}
		return nil, rest.ErrGiftPay.WithMessageAndCode(payResp.Code, payResp.Message, payResp.RequestId)
	}
	status := model.SendGiftStatusSuccess
	sResp := &SendGiftResponse{
		LiveId:   req.LiveId,
		UserId:   userId,
		BizId:    bizId,
		GiftId:   req.GiftId,
		Amount:   req.Amount,
		AnchorId: liveEntity.AnchorId,
		Status:   status,
	}
	err = s.UpdateGiftStatus(context, bizId, status)
	if err != nil {
		log.Errorf("update gift status error %s", err.Error())
		//该错误没有返回
	}

	notifyItem := BroadcastGiftNotifyItem{
		LiveId: liveEntity.LiveId,
		UserId: userId,
		GiftId: req.GiftId,
		Amount: req.Amount,
	}
	userInfo, err := user.GetService().FindUser(context, userId)
	if err != nil {
		log.Errorf("get user info for %s error %s", userId, err.Error())
		return sResp, nil
	}
	err = notify.SendNotifyToLive(context, userInfo, liveEntity, notify.ActionTypeGiftNotify, &notifyItem)
	if err != nil {
		log.Errorf("send notify to live %s error %s", liveEntity.LiveId, err.Error())
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
	RequestId string `json:"request_id"` //请求ID
	Code      int    `json:"code"`       //错误码，0 成功，其他失败
	Message   string `json:"message"`    //错误信息
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
		return rest.ErrInternal
	}
	return nil
}

func (s *ServiceImpl) UpdateGiftStatus(context context.Context, bizId string, status int) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err := db.Model(&model.LiveGift{}).Where("biz_id = ?", bizId).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceImpl) SearchGiftByAnchorId(context context.Context, anchorId string, pageNum, pageSize int) (liveGifts []*model.LiveGift, totalCount int, err error) {
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

func (s *ServiceImpl) SearchGiftByLiveId(context context.Context, liveId string, pageNum, pageSize int) (liveGifts []*model.LiveGift, totalCount int, err error) {
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

func (s *ServiceImpl) SearchGiftByUserId(context context.Context, userId string, pageNum, pageSize int) (liveGifts []*model.LiveGift, totalCount int, err error) {
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
