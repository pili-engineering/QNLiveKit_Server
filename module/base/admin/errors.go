package admin

import (
	"net/http"

	"github.com/qbox/livekit/core/rest"
)

var ErrUserPassword = &rest.Error{StatusCode: http.StatusOK, Code: 30001, Message: "Invalid username or password"}
