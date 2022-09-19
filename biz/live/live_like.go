package live

import (
	"context"
	"fmt"
	"time"

	"github.com/qbox/livekit/common/cache"
	"github.com/qbox/livekit/utils/logger"
)

// 点赞先在Redis 中进行缓存，周期性的将缓存数据刷到DB
// 缓存分几个部分
// 1. 直播间所有用户的点赞数，缓存时间 2 天，缓存结构 HASH
//      key: like:{live_id}
//      field: {user_id}， "_" 表示直播间总数
//
// 2. 缓存当前时间段内，有点赞的直播间，缓存时间 1 小时， 缓存结构 SET
//      key: like_rooms:{time}
//
// 3. 缓存当前时间段内，直播间内的点赞用户，缓存时间 1 小时， 缓存结构 SET
//      key: like_users:{live_id}:{time}

const (
	liveLikeFmt      = "like:%s"
	liveLikeRoomsFmt = "like_rooms:%d"
	liveLikeUsersFmt = "like_users:%s:%d"

	liveLikeTTL          = 48 * time.Hour
	liveLikeRoomsTTL     = time.Hour
	liveLikeRoomUsersTTL = time.Hour
)

func (s *Service) liveLikeKey(liveId string) string {
	return fmt.Sprintf(liveLikeFmt, liveId)
}

func (s *Service) liveLikeRoomsKey(now time.Time) string {
	return fmt.Sprintf(liveLikeRoomsFmt, now.Unix())
}

func (s *Service) liveLikeRoomUsersKey(now time.Time, liveId string) string {
	return fmt.Sprintf(liveLikeUsersFmt, liveId, now.Unix())
}

func (s *Service) cacheLike(ctx context.Context, liveId string, userId string, count int64) (my, total int64, err error) {
	log := logger.ReqLogger(ctx)
	my, total, err = s.incrRoomLikes(ctx, liveId, userId, count)
	if err != nil {
		log.Errorf("incrRoomLikes error %s", err.Error())
		return
	}

	now := time.Now()
	s.cacheLikeRooms(ctx, now, liveId)
	s.cacheLikeRoomUsers(ctx, now, liveId, userId)

	return
}

func (s *Service) incrRoomLikes(ctx context.Context, liveId string, userId string, count int64) (my, total int64, err error) {
	log := logger.ReqLogger(ctx)
	key := s.liveLikeKey(liveId)
	my, err = cache.Client.HIncrByCtx(ctx, key, userId, count)
	if err != nil {
		log.Errorf("cache error %s", err.Error())
		return 0, 0, err
	}
	cache.Client.Expire(key, liveLikeTTL)
	total, err = cache.Client.HIncrByCtx(ctx, key, "_", count)
	if err != nil {
		log.Errorf("cache error %s", err.Error())
		return 0, 0, err
	}

	return my, total, err
}

func (s *Service) cacheLikeRooms(ctx context.Context, now time.Time, liveId string) error {
	key := s.liveLikeRoomsKey(now)
	cache.Client.SAddCtx(ctx, key, []interface{}{liveId})
	cache.Client.Expire(key, liveLikeRoomsTTL)
	return nil
}

func (s *Service) getLikeRooms(ctx context.Context, now time.Time) ([]string, error) {
	key := s.liveLikeRoomsKey(now)
	return cache.Client.SMembersCtx(ctx, key)
}

func (s *Service) cacheLikeRoomUsers(ctx context.Context, now time.Time, liveId string, userId string) error {
	key := s.liveLikeRoomUsersKey(now, liveId)
	cache.Client.SAddCtx(ctx, key, []interface{}{userId})
	cache.Client.Expire(key, liveLikeRoomUsersTTL)
	return nil
}

func (s *Service) getLikeRoomUsers(ctx context.Context, now time.Time, liveId string) ([]string, error) {
	key := s.liveLikeRoomUsersKey(now, liveId)
	return cache.Client.SMembersCtx(ctx, key)
}
