package live

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/report"
	"github.com/qbox/livekit/common/cache"
	"github.com/qbox/livekit/common/mysql"
)

func TestService_cacheLikeRoomUsers(t *testing.T) {
	likeCacheSetup()
	defer likeCacheTearDown()

	ctx := context.Background()
	now := time.Now()
	liveId1 := fmt.Sprintf("live_%d_1", now.Unix())
	liveId2 := fmt.Sprintf("live_%d_2", now.Unix())

	userId1 := "user_1"
	userId2 := "user_2"
	userId3 := "user_3"

	s := &Service{}
	s.cacheLikeRoomUsers(ctx, now, liveId1, userId1)
	s.cacheLikeRoomUsers(ctx, now, liveId1, userId2)
	s.cacheLikeRoomUsers(ctx, now, liveId1, userId3)
	s.cacheLikeRoomUsers(ctx, now, liveId1, userId1)

	s.cacheLikeRoomUsers(ctx, now, liveId2, userId1)
	s.cacheLikeRoomUsers(ctx, now, liveId2, userId2)
	s.cacheLikeRoomUsers(ctx, now, liveId2, userId1)

	users, err := s.getLikeRoomUsers(ctx, now, liveId1)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(users))

	users, err = s.getLikeRoomUsers(ctx, now, liveId2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))
}

func TestService_cacheLikeRooms(t *testing.T) {
	likeCacheSetup()
	defer likeCacheTearDown()

	ctx := context.Background()
	now := time.Now()
	liveId1 := fmt.Sprintf("live_%d_1", now.Unix())
	liveId2 := fmt.Sprintf("live_%d_2", now.Unix())

	s := &Service{}
	err := s.cacheLikeRooms(ctx, now, liveId1)
	assert.Nil(t, err)
	err = s.cacheLikeRooms(ctx, now, liveId2)
	assert.Nil(t, err)

	rooms, err := s.getLikeRooms(ctx, now)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rooms))

	err = s.cacheLikeRooms(ctx, now, liveId2)
	assert.Nil(t, err)

	rooms, err = s.getLikeRooms(ctx, now)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rooms))
}

func TestService_incrRoomLikes(t *testing.T) {
	likeCacheSetup()
	defer likeCacheTearDown()

	ctx := context.Background()
	now := time.Now()
	liveId1 := fmt.Sprintf("live_%d_1", now.Unix())
	liveId2 := fmt.Sprintf("live_%d_2", now.Unix())

	userId1 := "user_1"
	userId2 := "user_2"
	userId3 := "user_3"

	s := &Service{}
	my, total, err := s.incrRoomLikes(ctx, liveId1, userId1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), my)
	assert.Equal(t, int64(1), total)

	my, total, _ = s.incrRoomLikes(ctx, liveId1, userId1, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), my)
	assert.Equal(t, int64(3), total)

	my, total, _ = s.incrRoomLikes(ctx, liveId1, userId2, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), my)
	assert.Equal(t, int64(4), total)

	my, total, _ = s.incrRoomLikes(ctx, liveId1, userId2, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), my)
	assert.Equal(t, int64(6), total)

	my, total, _ = s.incrRoomLikes(ctx, liveId2, userId2, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), my)
	assert.Equal(t, int64(1), total)

	my, total, _ = s.incrRoomLikes(ctx, liveId2, userId3, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), my)
	assert.Equal(t, int64(3), total)

}

func TestService_cacheLike(t *testing.T) {
	likeCacheSetup()
	defer likeCacheTearDown()

	ctx := context.Background()
	now := time.Now()
	liveId1 := fmt.Sprintf("live_%d_1", now.Unix())
	liveId2 := fmt.Sprintf("live_%d_2", now.Unix())

	userId1 := "user_1"
	userId2 := "user_2"
	userId3 := "user_3"

	s := &Service{}
	my, total, err := s.cacheLike(ctx, liveId1, userId1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), my)
	assert.Equal(t, int64(1), total)

	my, total, _ = s.cacheLike(ctx, liveId1, userId1, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), my)
	assert.Equal(t, int64(3), total)

	my, total, _ = s.cacheLike(ctx, liveId1, userId2, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), my)
	assert.Equal(t, int64(4), total)

	my, total, _ = s.cacheLike(ctx, liveId1, userId2, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), my)
	assert.Equal(t, int64(6), total)

	my, total, _ = s.cacheLike(ctx, liveId2, userId2, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), my)
	assert.Equal(t, int64(1), total)

	my, total, _ = s.cacheLike(ctx, liveId2, userId3, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), my)
	assert.Equal(t, int64(3), total)
}

func likeCacheSetup() {
	cache.Init(&cache.Config{
		Type:     cache.TypeNode,
		Addr:     "127.0.0.1:6379",
		Addrs:    nil,
		Password: "",
	})
}

func likeCacheTearDown() {

}

func TestService_getRoomLikes(t *testing.T) {
	likeCacheSetup()
	defer likeCacheTearDown()

	ctx := context.Background()
	now := time.Now()
	liveId1 := fmt.Sprintf("live_%d_1", now.Unix())

	userId1 := "user_1"
	userId2 := "user_2"
	userId3 := "user_3"

	s := &Service{}
	my, total, err := s.cacheLike(ctx, liveId1, userId1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), my)
	assert.Equal(t, int64(1), total)

	my, total, _ = s.cacheLike(ctx, liveId1, userId1, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), my)
	assert.Equal(t, int64(3), total)

	my, total, _ = s.cacheLike(ctx, liveId1, userId2, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), my)
	assert.Equal(t, int64(4), total)

	my, total, _ = s.cacheLike(ctx, liveId1, userId2, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), my)
	assert.Equal(t, int64(6), total)

	likeMap, err := s.getRoomLikes(ctx, liveId1, []string{userId1, userId2, userId3})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(likeMap))
	assert.Equal(t, int64(3), likeMap[userId1])
	assert.Equal(t, int64(3), likeMap[userId2])
}

func TestService_updateLastFlushTime(t *testing.T) {
	mysql.Init(&mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "live_test",
		Default:  "live",
	}, &mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "live_test",
		Default:  "live",
		ReadOnly: true,
	})
	mysql.GetLive().AutoMigrate(&model.LiveLikeFlush{})
	s := &Service{}
	tests := []struct {
		name     string
		lastTime int64
		now      int64
		want     int64
	}{
		{
			name:     "超前",
			lastTime: 5,
			now:      5,
			want:     5,
		},
		{
			name:     "未延迟",
			lastTime: 4,
			now:      5,
			want:     5,
		},
		{
			name:     "延迟小于5秒",
			lastTime: 4,
			now:      7,
			want:     7,
		},
		{
			name:     "延迟大于5秒",
			lastTime: 4,
			now:      13,
			want:     9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.updateLastFlushTime(context.Background(), tt.lastTime, tt.now)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_getLastFlushTime(t *testing.T) {
	mysql.Init(&mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "live_test",
		Default:  "live",
	}, &mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "live_test",
		Default:  "live",
		ReadOnly: true,
	})
	mysql.GetLive().AutoMigrate(&model.LiveLikeFlush{})

	s := &Service{}
	got, err := s.getLastFlushTime(context.Background())
	assert.Nil(t, err)
	assert.Greater(t, got, int64(0))
}

func TestService_getRoomLikeUsers(t *testing.T) {
	likeCacheSetup()
	defer likeCacheTearDown()

	ctx := context.Background()
	now := time.Now()
	liveId1 := fmt.Sprintf("live_%d_1", now.Unix())
	liveId2 := fmt.Sprintf("live_%d_2", now.Unix())

	userId1 := "user_1"
	userId2 := "user_2"
	userId3 := "user_3"

	s := &Service{}
	s.cacheLike(ctx, liveId1, userId1, 1)
	s.cacheLike(ctx, liveId1, userId2, 2)

	s.cacheLike(ctx, liveId2, userId1, 1)
	s.cacheLike(ctx, liveId2, userId2, 2)
	s.cacheLike(ctx, liveId2, userId3, 3)

	got, err := s.getRoomLikeUsers(context.Background(), now)
	assert.Nil(t, err)
	assert.NotEmpty(t, got)
}

func TestService_flushCacheLikes(t *testing.T) {
	mysql.Init(&mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "live",
		Default:  "live",
	}, &mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "live",
		Default:  "live",
		ReadOnly: true,
	})

	cache.Init(&cache.Config{
		Type:     cache.TypeNode,
		Addr:     "127.0.0.1:6379",
		Addrs:    nil,
		Password: "",
	})
	report.InitService()

	from := int64(1663720421)
	to := int64(1663720424)
	s := &Service{}
	s.flushCacheLikes(context.Background(), from, to)
}
