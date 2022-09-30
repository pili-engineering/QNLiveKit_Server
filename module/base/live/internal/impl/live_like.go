package impl

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/module/base/stats"
	"github.com/qbox/livekit/module/store/cache"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
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

func (s *Service) getRoomLikes(ctx context.Context, liveId string, userIds []string) (map[string]int64, error) {
	log := logger.ReqLogger(ctx)
	key := s.liveLikeKey(liveId)

	likes, err := cache.Client.HMGetCtx(ctx, key, userIds)
	if err != nil {
		log.Errorf("cache error %s", err.Error())
		return nil, err
	}

	ret := map[string]int64{}
	for i, _ := range likes {
		if l := likes[i]; l != nil {
			if lstr, ok := l.(string); ok {
				lint, _ := strconv.ParseInt(lstr, 10, 64)
				ret[userIds[i]] = lint
			}
		}
	}

	return ret, nil
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

// FlushCacheLikes 将缓存中的点赞数据，写入DB
// 为了避免服务器间的时钟不同步，这里重刷数秒
const flushSeconds = 5

func (s *Service) FlushCacheLikes(ctx context.Context) {
	log := logger.ReqLogger(ctx)
	for {
		now := time.Now().Unix()
		lastTime, err := s.getLastFlushTime(ctx)
		if err != nil {
			log.Errorf("getLastFlushTime error %s", err.Error())
			time.Sleep(time.Second)
			continue
		}
		if lastTime >= now {
			time.Sleep(time.Second)
			continue
		}

		from := lastTime - flushSeconds + 1
		to := lastTime + 1
		if err := s.flushCacheLikes(ctx, from, to); err != nil {
			log.Errorf("flushCacheLikes from %d to %d error %s", from, to, err.Error())
			time.Sleep(time.Second)
			continue
		}

		_, err = s.updateLastFlushTime(ctx, lastTime, now)
		if err != nil {
			log.Errorf("updateLastFlushTime error %s", err)
		}
	}
}

func (s *Service) updateLastFlushTime(ctx context.Context, lastTime, now int64) (int64, error) {
	if lastTime >= now {
		return lastTime, nil
	}

	if lastTime+flushSeconds >= now {
		lastTime = now
	} else {
		lastTime = lastTime + flushSeconds
	}

	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	flush := model.LiveLikeFlush{
		Id:             1,
		LastUpdateTime: lastTime,
	}
	err := db.Save(&flush).Error
	if err != nil {
		log.Errorf("save LiveLikeFlush error %s", err.Error())
	}

	return lastTime, err
}

func (s *Service) getLastFlushTime(ctx context.Context) (int64, error) {
	now := time.Now().Unix()
	flush := model.LiveLikeFlush{
		Id:             1,
		LastUpdateTime: now - 300,
	}
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	err := db.FirstOrCreate(&flush, "id = 1").Error
	if err != nil {
		log.Errorf("db.FirstOrCreate error %s", err.Error())
		return 0, err
	}
	return flush.LastUpdateTime, nil
}

func (s *Service) flushCacheLikes(ctx context.Context, from, to int64) error {
	log := logger.ReqLogger(ctx)
	roomUsers := map[string]map[string]struct{}{} //room -> user
	for t := from; t <= to; t++ {
		curRoomUsers, err := s.getRoomLikeUsers(ctx, time.Unix(t, 0))
		if err != nil {
			log.Errorf("getRoomLikeUsers error %s", err.Error())
			return err
		}

		if len(curRoomUsers) == 0 {
			continue
		}

		s.combineRoomUsers(roomUsers, curRoomUsers)
	}

	if len(roomUsers) == 0 {
		return nil
	}

	var err error = nil
	wg := sync.WaitGroup{}
	wg.Add(len(roomUsers))
	for room, users := range roomUsers {
		go func(room string, users map[string]struct{}) {
			defer wg.Done()

			err1 := s.flushRoomCacheLikes(ctx, room, users)
			if err1 != nil {
				log.Errorf("flushRoomCacheLikes for %s error %s", room, err.Error())
				err = err1
			}
		}(room, users)
	}
	wg.Wait()
	return err
}

func (s *Service) flushRoomCacheLikes(ctx context.Context, room string, users map[string]struct{}) error {
	log := logger.ReqLogger(ctx)
	userIds := make([]string, 0, len(users))
	for user, _ := range users {
		userIds = append(userIds, user)
	}

	likes, err := s.getRoomLikes(ctx, room, userIds)
	if err != nil {
		log.Errorf("getRoomLikes error %s", err.Error())
		return err
	}

	entities := make([]*model.StatsSingleLiveEntity, 0, len(likes))
	for userId, count := range likes {
		entities = append(entities, &model.StatsSingleLiveEntity{
			Type:      model.StatsTypeLike,
			LiveId:    room,
			UserId:    userId,
			BizId:     "",
			Count:     int(count),
			UpdatedAt: timestamp.Now(),
		})
	}
	err = stats.GetService().SaveStatsSingleLive(ctx, entities)
	if err != nil {
		log.Errorf("SaveStatsSingleLive error %s", err.Error())
	}
	return err
}

func (s *Service) combineRoomUsers(dst, src map[string]map[string]struct{}) {
	for room, users := range src {
		old := dst[room]
		if old == nil {
			old = map[string]struct{}{}
			dst[room] = old
		}

		for user, _ := range users {
			old[user] = struct{}{}
		}
	}
}

// getRoomLikeUsers 获取某一个秒内，对房间点赞的用户
func (s *Service) getRoomLikeUsers(ctx context.Context, now time.Time) (map[string]map[string]struct{}, error) {
	log := logger.ReqLogger(ctx)
	rooms, err := s.getLikeRooms(ctx, now)
	if err != nil {
		log.Errorf("getLikeRooms error %s", err.Error())
		return nil, err
	}

	if len(rooms) == 0 {
		return nil, nil
	}

	ret := map[string]map[string]struct{}{}
	for _, room := range rooms {
		users, err := s.getLikeRoomUsers(ctx, now, room)
		if err != nil {
			log.Errorf("getLikeRoomUsers error %s", err.Error())
			return nil, err
		}
		if len(users) == 0 {
			continue
		}

		roomUsers := ret[room]
		if roomUsers == nil {
			roomUsers = map[string]struct{}{}
			ret[room] = roomUsers
		}

		for _, user := range users {
			roomUsers[user] = struct{}{}
		}
	}

	return ret, nil
}
