package admin

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/notify"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/live"
	"github.com/qbox/livekit/module/base/user"
	"github.com/qbox/livekit/module/biz/censor/dto"
	"github.com/qbox/livekit/module/biz/censor/internal/impl"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

func RegisterRoutes() {
	//censorGroup := group.Group("/censor")
	//censorGroup.POST("/config", censorController.UpdateCensorConfig)
	//censorGroup.GET("/config", censorController.GetCensorConfig)
	//censorGroup.POST("/stoplive/:liveId", censorController.PostStopLive)
	//
	//censorGroup.POST("/job/start", censorController.CreateJob)
	//censorGroup.POST("/job/close", censorController.CloseJob)
	//censorGroup.GET("/job/list", censorController.ListAllJobs)
	//censorGroup.GET("/job/query", censorController.QueryJob)
	//
	//censorGroup.GET("/live", censorController.SearchCensorLive)
	//censorGroup.GET("/record", censorController.SearchRecordImage)
	//censorGroup.POST("/audit", censorController.AuditRecordImage)

	httpq.Handle(http.MethodPost, "/manager/censor/callback", censorController.CallbackCensorJob)

	httpq.AdminHandle(http.MethodPost, "/censor/config", censorController.UpdateCensorConfig)
	httpq.AdminHandle(http.MethodGet, "/censor/config", censorController.GetCensorConfig)
	httpq.AdminHandle(http.MethodPost, "/censor/stoplive/:liveId", censorController.PostStopLive)

	httpq.AdminHandle(http.MethodPost, "/censor/job/start", censorController.CreateJob)
	httpq.AdminHandle(http.MethodPost, "/censor/job/close", censorController.CloseJob)
	httpq.AdminHandle(http.MethodGet, "/censor/job/list", censorController.ListAllJobs)
	httpq.AdminHandle(http.MethodGet, "/censor/job/query", censorController.QueryJob)

	httpq.AdminHandle(http.MethodGet, "/censor/live", censorController.SearchCensorLive)
	httpq.AdminHandle(http.MethodGet, "/censor/record", censorController.SearchRecordImage)
	httpq.AdminHandle(http.MethodPost, "/censor/audit", censorController.AuditRecordImage)
}

var censorController = &CensorController{}

type CensorController struct {
}

// UpdateCensorConfig 更新配置
// return *dto.CensorConfigDto
func (c *CensorController) UpdateCensorConfig(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &dto.CensorConfigDto{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	if req.Enable && (req.Interval < 1 || req.Interval > 60) {
		log.Errorf("request interval invalid error")
		return nil, rest.ErrBadRequest
	}
	censorService := impl.GetInstance()
	err := censorService.UpdateCensorConfig(ctx, dto.CConfigDtoToEntity(req))
	if err != nil {
		log.Errorf(" UpdateCensorConfig error:%v", err)
		return nil, err
	}

	return req, nil
}

// GetCensorConfig 查询三鉴配置
// return *dto.CensorConfigDto
func (c *CensorController) GetCensorConfig(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	censorService := impl.GetInstance()
	censorConfig, err := censorService.GetCensorConfig(ctx)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	return dto.CConfigEntityToDto(censorConfig), nil
}

// PostStopLive 管理员强制停止直播
// return nil
func (c *CensorController) PostStopLive(ctx *gin.Context) (interface{}, error) {
	adminInfo := auth.GetAdminInfo(ctx)

	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	if liveId == "" {
		return nil, rest.ErrBadRequest.WithMessage("empty liveId")
	}

	liveEntity, err := live.GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("find live error %s", err.Error())
		return nil, err
	}

	anchorInfo, err := user.GetService().FindUser(ctx, liveEntity.AnchorId)
	if err != nil {
		log.Errorf("get anchor info for %s error %s", liveEntity.AnchorId, err.Error())
	}
	notifyItem := LiveNotifyItem{
		LiveId:  liveEntity.LiveId,
		Message: "直播涉嫌违规，\n管理员已关闭直播间。",
	}
	err = notify.SendNotifyToLive(ctx, anchorInfo, liveEntity, notify.ActionTypeCensorStop, &notifyItem)
	if err != nil {
		log.Errorf("send notify to live %s error %s", liveEntity.LiveId, err.Error())
	}

	err = live.GetService().AdminStopLive(ctx, liveId, model.LiveStopReasonCensor, adminInfo.UserId)
	if err != nil {
		log.Errorf("stop live failed, err: %v", err)
		return nil, err
	}

	return nil, nil
}

// CallbackCensorJob 处理回调
// return
func (c *CensorController) CallbackCensorJob(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &CensorCallBack{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	if req.Image.Code != 200 {
		log.Errorf("CallbackCensorJob  response Error %v", req.Error.Message)
		return nil, fmt.Errorf("image.code %d", req.Image.Code)
	}
	url := impl.GetInstance().ImageBucketToUrl(req.Image.Url)
	m := &model.CensorImage{
		JobID:      req.Image.Job,
		Url:        url,
		CreatedAt:  req.Image.Timestamp,
		Suggestion: req.Image.Result.Suggestion,
		Pulp:       req.Image.Result.Scenes.Pulp.Suggestion,
		Ads:        req.Image.Result.Scenes.Ads.Suggestion,
		Politician: req.Image.Result.Scenes.Politician.Suggestion,
		Terror:     req.Image.Result.Scenes.Terror.Suggestion,
	}

	censorService := impl.GetInstance()
	liveCensor, err := censorService.GetLiveCensorJobByJobId(ctx, req.Image.Job)
	if err != nil {
		log.Errorf("GetLiveCensorJobByJobId error %v", err)
		return nil, err
	}
	m.LiveID = liveCensor.LiveID
	err = censorService.SaveCensorImage(ctx, m)
	if err != nil {
		log.Errorf("SaveCensorImage error %v", err)
		return nil, err
	}
	err = live.GetService().UpdateLiveRelatedReview(ctx, m.LiveID, &req.Image.Timestamp)
	if err != nil {
		log.Errorf("UpdateLiveRelatedReview error %v", err)
		return nil, err
	}

	return nil, nil
}

// CreateJob 创建三鉴任务
// return nil
func (c *CensorController) CreateJob(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &CensorCreateRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	liveEntity, err := live.GetService().LiveInfo(ctx, req.LiveId)
	if err != nil {
		log.Errorf("LiveInfo error:%v", err)
		return nil, err
	}
	censorService := impl.GetInstance()
	err = censorService.CreateCensorJob(ctx, liveEntity)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		return nil, err
	}
	return nil, nil
}

// CloseJob 关闭三鉴任务
func (c *CensorController) CloseJob(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &CensorCloseRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	censorService := impl.GetInstance()
	err := censorService.StopCensorJob(ctx, req.LiveId)
	if err != nil {
		log.Errorf("Stop censor job  error:%v", err)
		return nil, err
	}
	return nil, nil
}

type SearchRecordRequest struct {
	PageNum  int     `json:"page_num" form:"page_num"`
	PageSize int     `json:"page_size" form:"page_size"`
	LiveId   *string `json:"live_id" form:"live_id"`
	IsReview *int    `json:"is_review" form:"is_review"`
}

// SearchRecordImage 查询三鉴图片结果
// return
func (c *CensorController) SearchRecordImage(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &SearchRecordRequest{}
	if err := ctx.BindQuery(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if req.IsReview != nil {
		if *req.IsReview != impl.IsReviewNo && *req.IsReview != impl.IsReviewYes {
			log.Errorf(" invalid argument")
			return nil, rest.ErrBadRequest
		}
	}
	images, count, err := impl.GetInstance().SearchCensorImage(ctx, req.IsReview, req.PageNum, req.PageSize, req.LiveId)
	if err != nil {
		log.Errorf("search censor image  failed, err: %v", err)
		return nil, err
	}

	endPage := false
	if len(images) < req.PageSize {
		endPage = true
	}

	imageDtos := make([]*dto.CensorImageDto, len(images))
	for i, image := range images {
		imageDtos[i] = dto.CensorImageModelToDto(&image)
	}

	pageResult := &rest.PageResult{
		TotalCount: count,
		PageTotal:  int(math.Ceil(float64(count) / float64(req.PageSize))),
		EndPage:    endPage,
		List:       imageDtos,
	}

	return pageResult, nil
}

type SearchCensorLiveRequest struct {
	PageNum  int  `json:"page_num" form:"page_num"`
	PageSize int  `json:"page_size" form:"page_size"`
	IsReview *int `json:"is_review" form:"is_review"`
}

// SearchCensorLive 查询待审核直播间
// return *rest.PageResult
func (c *CensorController) SearchCensorLive(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &SearchCensorLiveRequest{}
	if err := ctx.BindQuery(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if req.IsReview != nil {
		if *req.IsReview != impl.IsReviewNo {
			log.Errorf(" invalid argument, is_review %d", *req.IsReview)
			return nil, rest.ErrBadRequest
		}
	}

	lives, count, err := impl.GetInstance().SearchCensorLive(ctx, req.IsReview, req.PageNum, req.PageSize)
	if err != nil {
		log.Errorf("search censor live  failed, err: %v", err)
		return nil, err
	}

	for i, liveEntity := range lives {
		anchor, err := live.GetService().FindLiveRoomUser(ctx, liveEntity.LiveId, liveEntity.AnchorId)
		if err != nil {
			if !errors.Is(err, rest.ErrNotFound) {
				log.Errorf("FindLiveRoomUser failed, err: %v", err)
				return nil, rest.ErrInternal
			}
		} else {
			lives[i].AnchorStatus = int(anchor.Status)
		}
		anchor2, err := user.GetService().FindUser(ctx, liveEntity.AnchorId)
		if err != nil {
			log.Errorf("FindUser  failed, err: %v", err)
			return nil, err
		}
		lives[i].Nick = anchor2.Nick
		lives[i].AnchorStatus = int(anchor.Status)
	}

	endPage := false
	if len(lives) < req.PageSize {
		endPage = true
	}

	pageResult := &rest.PageResult{
		TotalCount: count,
		PageTotal:  int(math.Ceil(float64(count) / float64(req.PageSize))),
		EndPage:    endPage,
		List:       lives,
	}
	return pageResult, nil
}

// ListAllJobs 查看所有的三鉴任务
// return
func (c *CensorController) ListAllJobs(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &CensorListRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	request := &impl.JobListRequest{
		Start:  req.Start.Unix(),
		End:    req.End.Unix(),
		Status: req.Status,
		Limit:  req.Limit,
		Marker: req.Marker,
	}
	resp := &impl.JobListResponse{}
	err := impl.GetInstance().Client.JobList(ctx, request, resp)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		return nil, err
	}

	return resp.Data, nil
}

// QueryJob 查询三鉴任务
// return
func (c *CensorController) QueryJob(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &CensorQueryRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	requset := &impl.JobQueryRequest{
		Start:       req.Start.Unix(),
		End:         req.End.Unix(),
		Job:         req.Job,
		Suggestions: req.Suggestions,
	}
	resp := &impl.JobQueryResponse{}
	err := impl.GetInstance().Client.JobQuery(ctx, requset, resp)
	if err != nil {
		log.Errorf("Job Query error:%v", err)
		return nil, err
	}
	return resp.Data, nil
}

// AuditRecordImage 管理员审核三鉴结果
// return nil
func (c *CensorController) AuditRecordImage(ctx *gin.Context) (interface{}, error) {
	userInfo := auth.GetAdminInfo(ctx)
	log := logger.ReqLogger(ctx)
	req := &CensorAuditRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	if req.ReviewAnswer != model.AuditResultPass && req.ReviewAnswer != model.AuditResultBlock {
		log.Errorf("invalid request %+v", req)
		return nil, rest.ErrBadRequest
	}
	if len(req.Images) == 0 {
		log.Errorf("invalid request %+v", req)
		return nil, rest.ErrBadRequest
	}
	censorService := impl.GetInstance()
	updates := map[string]interface{}{}
	updates["is_review"] = 1
	updates["review_answer"] = req.ReviewAnswer
	updates["review_time"] = timestamp.Now()
	updates["review_user_id"] = userInfo.UserId
	err := censorService.BatchUpdateCensorImage(ctx, req.Images, updates)
	if err != nil {
		log.Errorf("update audit censor image info  error %s", err.Error())
		return nil, err
	}
	err = live.GetService().UpdateLiveRelatedReview(ctx, req.LiveId, nil)
	if err != nil {
		log.Errorf("update Live Related Review error %s", err.Error())
		return nil, err
	}

	if req.Notify {
		if req.ReviewAnswer == model.AuditResultBlock {
			go c.notifyCensorBlock(ctx, req.LiveId)
		}
	}
	return nil, nil
}

func (c *CensorController) notifyCensorBlock(ctx context.Context, liveId string) {
	log := logger.ReqLogger(ctx)
	anchor, err := live.GetService().GetLiveAuthor(ctx, liveId)
	if err != nil {
		log.Errorf("get live %s error %v", liveId, err)
		return
	}

	notifyItem := LiveNotifyItem{
		LiveId:  liveId,
		Message: "请注意您的直播内容\\n如严重违规，管理员将强行关闭直播间。",
	}
	err = notify.SendNotifyToUser(ctx, anchor, notify.ActionTypeCensorNotify, &notifyItem)
	if err != nil {
		log.Errorf("send notify to user error %v", err)
	}
}

type LiveNotifyItem struct {
	LiveId  string `json:"live_id"`
	Message string `json:"message"`
}

type CensorCreateRequest struct {
	LiveId string `json:"live_id"`
}

type CensorAuditRequest struct {
	LiveId       string `json:"live_id"`
	Images       []uint `json:"image_list"`
	Notify       bool   `json:"notify"` //是否发送违规警告
	ReviewAnswer int    `json:"review_answer"`
}

type CensorListRequest struct {
	Start  timestamp.Timestamp `json:"start"`
	End    timestamp.Timestamp `json:"end"`
	Status string              `json:"status"`
	Limit  int                 `json:"limit"`
	Marker string              `json:"marker"`
}

type CensorCloseRequest struct {
	LiveId string `json:"live_id"`
}

type CensorQueryRequest struct {
	Start       timestamp.Timestamp `json:"start"`
	End         timestamp.Timestamp `json:"end"`
	Job         string              `json:"job"`
	Suggestions []string            `json:"suggestions"`
}

type CensorCallBack struct {
	Job  string `json:"job"`
	Live struct {
		Id   string `json:"id"`
		Uri  string `json:"uri"`
		Info string `json:"info"`
	} `json:"live"`
	Error struct {
		Timestamp int    `json:"timestamp"`
		Message   string `json:"message"`
	} `json:"error"`
	Image struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		Job       string `json:"job"`
		Url       string `json:"url"`
		Timestamp int    `json:"timestamp"`
		Result    struct {
			Suggestion string `json:"suggestion"`
			Scenes     struct {
				Ads struct {
					Suggestion string `json:"suggestion"`
					Details    []struct {
						Suggestion string  `json:"suggestion"`
						Label      string  `json:"label"`
						Score      float64 `json:"score"`
						Review     bool    `json:"review"`
					} `json:"details"`
				} `json:"ads"`
				Politician struct {
					Suggestion string `json:"suggestion"`
					Details    []struct {
						Suggestion string  `json:"suggestion"`
						Label      string  `json:"label"`
						Score      float64 `json:"score"`
						Review     bool    `json:"review"`
					} `json:"details"`
				} `json:"politician"`
				Pulp struct {
					Suggestion string `json:"suggestion"`
					Details    []struct {
						Suggestion string  `json:"suggestion"`
						Label      string  `json:"label"`
						Score      float64 `json:"score"`
						Review     bool    `json:"review"`
					} `json:"details"`
				} `json:"pulp"`
				Terror struct {
					Suggestion string `json:"suggestion"`
					Details    []struct {
						Suggestion string  `json:"suggestion"`
						Label      string  `json:"label"`
						Score      float64 `json:"score"`
						Review     bool    `json:"review"`
					} `json:"details"`
				} `json:"terror"`
			} `json:"scenes"`
		} `json:"result"`
	} `json:"image"`
	Audio struct {
		Result struct {
			Suggestion string `json:"suggestion"`
		} `json:"result"`
	} `json:"audio"`
}
