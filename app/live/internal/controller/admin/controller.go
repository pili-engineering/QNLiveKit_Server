package admin

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/controller/server"
	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/admin"
	"github.com/qbox/livekit/biz/live"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/notify"
	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/biz/user"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

func RegisterCensorRoutes(group *gin.RouterGroup) {
	censorGroup := group.Group("/censor")
	censorGroup.POST("/config", censorController.UpdateCensorConfig)
	censorGroup.GET("/config", censorController.GetCensorConfig)
	censorGroup.POST("/stoplive/:liveId", censorController.PostStopLive)

	censorGroup.POST("/job/start", censorController.CreateJob)
	censorGroup.POST("/job/close", censorController.CloseJob)
	censorGroup.GET("/job/list", censorController.ListAllJobs)
	censorGroup.GET("/job/query", censorController.QueryJob)

	censorGroup.GET("/live", censorController.SearchCensorLive)
	censorGroup.GET("/record", censorController.SearchRecordImage)
	censorGroup.POST("/audit", censorController.AuditRecordImage)
}

var censorController = &CensorController{}

type CensorController struct {
}

type LoginRequest struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (c *CensorController) LoginManager(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &LoginRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	manService := admin.GetManagerService()
	admin, err := manService.FindAdminByUserName(ctx, req.UserName)
	if err != nil {
		log.Errorf("userName:%s, login error:%v", req.UserName, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	} else if admin.Password != req.Password {
		log.Errorf("userName:%s, login error:%v", req.UserName, api.ErrorLoginWrong)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), api.ErrorLoginWrong))
		return
	}

	authToken := token.AuthToken{
		UserId: admin.UserId,
		Role:   "admin",
	}

	tokenService := token.GetService()
	if token, err := tokenService.GenAuthToken(&authToken); err != nil {
		log.Errorf("")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	} else {
		resp := &server.GetAuthTokenResponse{
			Response: api.Response{
				RequestId: log.ReqID(),
				Code:      0,
				Message:   "success",
			},
		}
		resp.Data.AccessToken = token
		resp.Data.ExpiresAt = authToken.ExpiresAt
		ctx.JSON(http.StatusOK, resp)
	}
}

type CensorConfigResponse struct {
	api.Response
	Data *dto.CensorConfigDto `json:"data"`
}

func (c *CensorController) UpdateCensorConfig(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &dto.CensorConfigDto{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	censorService := admin.GetCensorService()
	err := censorService.UpdateCensorConfig(ctx, dto.CConfigDtoToEntity(req))
	if err != nil {
		log.Errorf(" UpdateCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &CensorConfigResponse{
		Response: api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
		Data: req,
	})
}

func (c *CensorController) GetCensorConfig(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	censorService := admin.GetCensorService()
	censorConfig, err := censorService.GetCensorConfig(ctx)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &CensorConfigResponse{
		Response: api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
		Data: dto.CConfigEntityToDto(censorConfig),
	})
}

func (c *CensorController) PostStopLive(ctx *gin.Context) {
	adminInfo := liveauth.GetAdminInfo(ctx)

	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	if liveId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	liveEntity, err := live.GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("find live error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
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
		ctx.JSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.Response{
		Code:      200,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

func (c *CensorController) CallbackCensorJob(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorCallBack{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	if req.Image.Code != 200 {
		log.Errorf("CallbackCensorJob  response Error %v", req.Error.Message)
		return
	}
	m := &model.CensorImage{
		JobID:      req.Image.Job,
		Url:        req.Image.Url,
		CreatedAt:  req.Image.Timestamp,
		Suggestion: req.Image.Result.Suggestion,
		Pulp:       req.Image.Result.Scenes.Pulp.Suggestion,
		Ads:        req.Image.Result.Scenes.Ads.Suggestion,
		Politician: req.Image.Result.Scenes.Politician.Suggestion,
		Terror:     req.Image.Result.Scenes.Terror.Suggestion,
	}

	censorService := admin.GetCensorService()
	liveCensor, err := censorService.GetLiveCensorJobByJobId(ctx, req.Image.Job)
	if err != nil {
		log.Errorf("GetLiveCensorJobByJobId error %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	m.LiveID = liveCensor.LiveID
	err = censorService.SaveCensorImage(ctx, m)
	if err != nil {
		log.Errorf("SaveCensorImage error %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	err = live.GetService().UpdateLiveRelatedReview(ctx, m.LiveID, &req.Image.Timestamp)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, api.Response{
		Code:      http.StatusOK,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

func (c *CensorController) CreateJob(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorCreateRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	liveEntity, err := live.GetService().LiveInfo(ctx, req.LiveId)
	if err != nil {
		log.Errorf("LiveInfo error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	censorService := admin.GetCensorService()
	err = censorService.CreateCensorJob(ctx, liveEntity)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &api.Response{
		RequestId: log.ReqID(),
		Code:      0,
		Message:   "success",
	})

}

func (c *CensorController) CloseJob(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorCloseRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	censorService := admin.GetCensorService()
	err := censorService.StopCensorJob(ctx, req.LiveId)
	if err != nil {
		log.Errorf("Stop censor job  error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &api.Response{
		RequestId: log.ReqID(),
		Code:      0,
		Message:   "success",
	})
}

func (c *CensorController) SearchRecordImage(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	audit := ctx.DefaultQuery("audit", strconv.Itoa(admin.AuditAll))
	pageNum := ctx.DefaultQuery("page_num", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	liveId := ctx.Query("live_id")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_size is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	auditInt, err := strconv.Atoi(audit)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "is_review is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	images, count, err := admin.GetCensorService().SearchCensorImage(ctx, auditInt, pageNumInt, pageSizeInt, liveId)
	if err != nil {
		log.Errorf("search censor image  failed, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "search  censor image failed",
			RequestId: log.ReqID(),
		})
		return
	}

	endPage := false
	if len(images) < pageSizeInt {
		endPage = true
	}
	response := &CensorImageListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = count
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	response.Data.List = images
	ctx.JSON(http.StatusOK, response)
}

func (c *CensorController) SearchCensorLive(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	audit := ctx.DefaultQuery("audit", strconv.Itoa(admin.AuditAll))
	pageNum := ctx.DefaultQuery("page_num", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_size is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	auditInt, err := strconv.Atoi(audit)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "is_review is not int",
			RequestId: log.ReqID(),
		})
		return
	}

	if auditInt != admin.AuditAll && auditInt != admin.AuditNo {
		log.Errorf(" invalid argument %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	lives, count, err := admin.GetCensorService().SearchCensorLive(ctx, auditInt, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("search censor live  failed, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "search  censor live failed",
			RequestId: log.ReqID(),
		})
		return
	}

	for i, liveEntity := range lives {
		anchor, err := live.GetService().FindLiveRoomUser(ctx, liveEntity.LiveId, liveEntity.AnchorId)
		if err != nil {
			log.Errorf("FindLiveRoomUser failed, err: %v", err)
			ctx.JSON(http.StatusInternalServerError, api.Response{
				Code:      http.StatusInternalServerError,
				Message:   "FindLiveRoomUser failed",
				RequestId: log.ReqID(),
			})
			return
		}
		anchor2, err := user.GetService().FindUser(ctx, liveEntity.AnchorId)
		if err != nil {
			log.Errorf("FindUser  failed, err: %v", err)
			ctx.JSON(http.StatusInternalServerError, api.Response{
				Code:      http.StatusInternalServerError,
				Message:   "FindUser failed",
				RequestId: log.ReqID(),
			})
			return
		}
		lives[i].Nick = anchor2.Nick
		lives[i].AnchorStatus = int(anchor.Status)
	}

	endPage := false
	if len(lives) < pageSizeInt {
		endPage = true
	}
	response := &CensorLiveListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = count
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	response.Data.List = lives
	ctx.JSON(http.StatusOK, response)
}

func (c *CensorController) ListAllJobs(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorListRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	request := &admin.JobListRequest{
		Start:  req.Start.Unix(),
		End:    req.End.Unix(),
		Status: req.Status,
		Limit:  req.Limit,
		Marker: req.Marker,
	}
	resp := &admin.JobListResponse{}
	err := admin.GetJobService().JobList(ctx, request, resp)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *CensorController) QueryJob(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &CensorQueryRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	requset := &admin.JobQueryRequest{
		Start:       req.Start.Unix(),
		End:         req.End.Unix(),
		Job:         req.Job,
		Suggestions: req.Suggestions,
	}
	resp := &admin.JobQueryResponse{}
	err := admin.GetJobService().JobQuery(ctx, requset, resp)
	if err != nil {
		log.Errorf("Job Query error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *CensorController) AuditRecordImage(ctx *gin.Context) {
	userInfo := ctx.MustGet(liveauth.AdminCtxKey).(*liveauth.AdminInfo)
	log := logger.ReqLogger(ctx)
	req := &CensorAuditRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	if req.ReviewAnswer != model.AuditResultPass && req.ReviewAnswer != model.AuditResultBlock {
		log.Errorf("invalid request %+v", req)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	if len(req.Images) == 0 {
		log.Errorf("invalid request %+v", req)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	censorService := admin.GetCensorService()
	updates := map[string]interface{}{}
	updates["is_review"] = 1
	updates["review_answer"] = req.ReviewAnswer
	updates["review_time"] = timestamp.Now()
	updates["review_user_id"] = userInfo.UserId
	err := censorService.BatchUpdateCensorImage(ctx, req.Images, updates)
	if err != nil {
		log.Errorf("update audit censor image info  error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	err = live.GetService().UpdateLiveRelatedReview(ctx, req.LiveId, nil)
	if err != nil {
		log.Errorf("update Live Related Review error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	if req.ReviewAnswer == model.AuditResultBlock {
		go c.notifyCensorBlock(ctx, req.LiveId)
	}

	ctx.JSON(http.StatusOK, api.Response{
		Code:      http.StatusOK,
		Message:   "success",
		RequestId: log.ReqID(),
	})
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
	ReviewAnswer int    `json:"audit_answer"`
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

type CensorImageListResponse struct {
	api.Response
	Data struct {
		TotalCount int                 `json:"total_count"`
		PageTotal  int                 `json:"page_total"`
		EndPage    bool                `json:"end_page"`
		List       []model.CensorImage `json:"list"`
	} `json:"data"`
}

type CensorLiveListResponse struct {
	api.Response
	Data struct {
		TotalCount int                `json:"total_count"`
		PageTotal  int                `json:"page_total"`
		EndPage    bool               `json:"end_page"`
		List       []admin.CensorLive `json:"list"`
	} `json:"data"`
}
