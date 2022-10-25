package auth

import (
	"context"

	"github.com/qbox/livekit/utils/logger"
)

func GetAdminInfo(ctx context.Context) *AdminInfo {
	log := logger.ReqLogger(ctx)
	i := ctx.Value(AdminCtxKey)
	if i == nil {
		return nil
	}

	if t, ok := i.(*AdminInfo); ok {
		return t
	} else {
		log.Errorf("%+v not user info", i)
		return nil
	}
}

func GetUserInfo(ctx context.Context) *UserInfo {
	log := logger.ReqLogger(ctx)
	i := ctx.Value(UserCtxKey)
	if i == nil {
		return nil
	}

	if t, ok := i.(*UserInfo); ok {
		return t
	} else {
		log.Errorf("%+v not user info", i)
		return nil
	}
}
