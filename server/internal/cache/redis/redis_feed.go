package redis

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_architect_social/server/internal/cache"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

var _ cache.Feed = (*FeedRedisCache)(nil)

const cacheUserKey = "feed:user_id:%d"

type FeedRedisCache struct {
	pool   *redis.Pool
	maxLen int
}

func NewFeedCache(pool *redis.Pool, maxLen int) *FeedRedisCache {
	return &FeedRedisCache{
		pool:   pool,
		maxLen: maxLen,
	}
}

func (r *FeedRedisCache) Add(userID int, f *repository.Feed) error {
	log.Debugf("caching message for userId: %d", userID)
	rc := r.pool.Get()
	defer rc.Close()

	fByte, err := json.Marshal(f)
	if err != nil {
		return fmt.Errorf("marshal feed error: %v", err)
	}

	k := fmt.Sprintf(cacheUserKey, userID)
	err = rc.Send("MULTI")
	if err != nil {
		return fmt.Errorf("marshal feed error: %v", err)
	}

	err = rc.Send("LPUSH", k, fByte)
	if err != nil {
		return fmt.Errorf("marshal feed error: %v", err)
	}

	err = rc.Send("LTRIM", k, 0, (r.maxLen - 1))
	if err != nil {
		return fmt.Errorf("marshal feed error: %v", err)
	}

	_, err = rc.Do("EXEC")
	if err != nil {
		return fmt.Errorf("add cache error: %v", err)
	}
	log.Debugf("cached message for userId: %d, with key: %s", userID, k)

	return nil
}

func (r *FeedRedisCache) GetUserFeed(userID int) ([]*repository.Feed, error) {
	rc := r.pool.Get()
	defer rc.Close()

	k := fmt.Sprintf(cacheUserKey, userID)
	log.Debugf("getting cache for userId: %d, with key: %s", userID, k)
	v, err := redis.Values(rc.Do("LRANGE", k, 0, -1))
	if err != nil {
		return nil, err
	}

	result := make([]*repository.Feed, 0, len(v))

	vals, err := redis.ByteSlices(v, nil)
	if err != nil {
		return nil, err
	}

	for _, vByte := range vals {
		var f *repository.Feed

		err := json.Unmarshal(vByte, &f)
		if err != nil {
			log.Errorf("unmarshal msg from redis error: %v", err)

			continue
		}

		result = append(result, f)
	}

	return result, nil
}
