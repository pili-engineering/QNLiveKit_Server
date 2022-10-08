package impl

import (
	"net/http"

	"github.com/qbox/livekit/core/rest"
)

var ErrGiftPay = &rest.Error{StatusCode: http.StatusAccepted, Code: 20001, Message: "PayGift Fail Biz Server error"}
