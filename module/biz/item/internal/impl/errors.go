package impl

import (
	"net/http"

	"github.com/qbox/livekit/core/rest"
)

var ErrLiveItemExceed = &rest.Error{StatusCode: http.StatusBadRequest, Code: 20001, Message: "items exceed in live room"}
