package live

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/qbox/livekit/common/cache"
)

func TestService_cacheLikeRoomUsers(t *testing.T) {
	type args struct {
		ctx    context.Context
		now    time.Time
		liveId string
		userId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			tt.wantErr(t, s.cacheLikeRoomUsers(tt.args.ctx, tt.args.now, tt.args.liveId, tt.args.userId), fmt.Sprintf("cacheLikeRoomUsers(%v, %v, %v, %v)", tt.args.ctx, tt.args.now, tt.args.liveId, tt.args.userId))
		})
	}
}

func TestService_cacheLikeRooms(t *testing.T) {
	type args struct {
		ctx    context.Context
		now    time.Time
		liveId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			tt.wantErr(t, s.cacheLikeRooms(tt.args.ctx, tt.args.now, tt.args.liveId), fmt.Sprintf("cacheLikeRooms(%v, %v, %v)", tt.args.ctx, tt.args.now, tt.args.liveId))
		})
	}
}

func TestService_incrRoomLikes(t *testing.T) {
	type args struct {
		ctx    context.Context
		liveId string
		userId string
		count  int64
	}
	tests := []struct {
		name      string
		args      args
		wantMy    int64
		wantTotal int64
		wantErr   assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			gotMy, gotTotal, err := s.incrRoomLikes(tt.args.ctx, tt.args.liveId, tt.args.userId, tt.args.count)
			if !tt.wantErr(t, err, fmt.Sprintf("incrRoomLikes(%v, %v, %v, %v)", tt.args.ctx, tt.args.liveId, tt.args.userId, tt.args.count)) {
				return
			}
			assert.Equalf(t, tt.wantMy, gotMy, "incrRoomLikes(%v, %v, %v, %v)", tt.args.ctx, tt.args.liveId, tt.args.userId, tt.args.count)
			assert.Equalf(t, tt.wantTotal, gotTotal, "incrRoomLikes(%v, %v, %v, %v)", tt.args.ctx, tt.args.liveId, tt.args.userId, tt.args.count)
		})
	}
}

func TestService_cacheLike(t *testing.T) {
	type args struct {
		ctx    context.Context
		liveId string
		userId string
		count  int64
	}
	tests := []struct {
		name      string
		args      args
		wantMy    int64
		wantTotal int64
		wantErr   assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			gotMy, gotTotal, err := s.cacheLike(tt.args.ctx, tt.args.liveId, tt.args.userId, tt.args.count)
			if !tt.wantErr(t, err, fmt.Sprintf("cacheLike(%v, %v, %v, %v)", tt.args.ctx, tt.args.liveId, tt.args.userId, tt.args.count)) {
				return
			}
			assert.Equalf(t, tt.wantMy, gotMy, "cacheLike(%v, %v, %v, %v)", tt.args.ctx, tt.args.liveId, tt.args.userId, tt.args.count)
			assert.Equalf(t, tt.wantTotal, gotTotal, "cacheLike(%v, %v, %v, %v)", tt.args.ctx, tt.args.liveId, tt.args.userId, tt.args.count)
		})
	}
}

func likeCacheSetup() {
	cache.Init(&cache.Config{
		Type:     cache.TypeNode,
		Addr:     "",
		Addrs:    nil,
		Password: "",
	})
}

func likeCacheTearDown() {

}
