package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Nil not found err
var Nil = redis.Nil

// Cache ...
type Cache struct {
	client redis.Cmdable
}

// Client ...
var Client *Cache

// Init  initialize cache client
func Init(conf *Config) error {
	var client redis.Cmdable
	switch conf.Type {
	case TypeCluster:
		client = getClusterClient(conf)
	case TypeNode:
		client = getClient(conf)
	default:
		return fmt.Errorf("invalid type %s", conf.Type)
	}
	Client = &Cache{
		client: client,
	}

	return nil
}

func getClient(conf *Config) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		MaxRetries:   maxRetries,
		MinIdleConns: idleConns,
		PoolSize:     poolSize,
	})
}

func getClusterClient(conf *Config) redis.Cmdable {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        conf.Addrs,
		Password:     conf.Password,
		MaxRetries:   maxRetries,
		MinIdleConns: idleConns,
		PoolSize:     poolSize,
	})
}

// Get get a key from redis cache
func (c *Cache) Get(key string) (string, error) {
	return c.client.Get(context.Background(), key).Result()
}

// GetInt64 get int64 key
func (c *Cache) GetInt64(key string) (int64, error) {
	i, err := c.client.Get(context.Background(), key).Int64()
	if err != nil && err != Nil {
		return i, err
	}

	return i, nil
}

// GetInt64OK return int and bool (false: if key not exist)
func (c *Cache) GetInt64OK(key string) (int64, bool) {
	i, err := c.client.Get(context.Background(), key).Int64()
	if err == Nil {
		return i, false
	}

	if err != nil {
		return i, false
	}

	return i, true

}

// Set a key val pair to redis cache
func (c *Cache) Set(key, val string, expiration time.Duration) error {
	return c.client.Set(context.Background(), key, val, expiration).Err()
}

// SetNX only if key donot exist then set key to value
func (c *Cache) SetNX(key, val string, expiration time.Duration) (bool, error) {
	ok, err := c.client.SetNX(context.Background(), key, val, expiration).Result()
	return ok, err
}

// SetEX set a key to value with expiration
func (c *Cache) SetEX(key, val string, expiration time.Duration) error {
	return c.client.Set(context.Background(), key, val, expiration).Err()
}

func (c *Cache) Expire(key string, expiration time.Duration) error {
	return c.client.Expire(context.Background(), key, expiration).Err()
}

// Del delete a key from cache
func (c *Cache) Del(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

// Increase increase a key atomicly
func (c *Cache) Incr(key string) error {
	return c.client.Incr(context.Background(), key).Err()
}

func (c *Cache) IncrBy(key string, value int64) (int64, error) {
	ret := c.client.IncrBy(context.Background(), key, value)
	return ret.Result()
}

// Exist check if a key is existed in cache
func (c *Cache) Exist(key string) bool {
	return c.client.Exists(context.Background(), key).Val() == 1
}

// EvalSha 执行load 到内存里面的 lua 脚本
func (c *Cache) EvalSha(sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return c.client.EvalSha(context.Background(), sha1, keys, args...).Result()
}

// Eval 执行 lua 脚本
func (c *Cache) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return c.client.Eval(context.Background(), script, keys, args...).Result()
}

// ScriptLoad load lua 脚本
func (c *Cache) ScriptLoad(script string) (string, error) {
	return c.client.ScriptLoad(context.Background(), script).Result()
}

func (c *Cache) HIncrBy(key, field string, incr int64) (int64, error) {
	return c.HIncrByCtx(context.Background(), key, field, incr)
}

func (c *Cache) HIncrByCtx(ctx context.Context, key, field string, incr int64) (int64, error) {
	return c.client.HIncrBy(ctx, key, field, incr).Result()
}

func (c *Cache) HMGet(key string, fields []string) ([]interface{}, error) {
	return c.HMGetCtx(context.Background(), key, fields)
}

func (c *Cache) HMGetCtx(ctx context.Context, key string, fields []string) ([]interface{}, error) {
	return c.client.HMGet(ctx, key, fields...).Result()
}

func (c *Cache) SAdd(key string, members []interface{}) (int64, error) {
	return c.SAddCtx(context.Background(), key, members)
}

func (c *Cache) SAddCtx(ctx context.Context, key string, members []interface{}) (int64, error) {
	return c.client.SAdd(ctx, key, members...).Result()
}

func (c *Cache) SMembers(key string) ([]string, error) {
	return c.SMembersCtx(context.Background(), key)
}

func (c *Cache) SMembersCtx(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

func (c *Cache) HGetAll(key string) (result map[string]string, err error) {
	return c.client.HGetAll(context.Background(), key).Result()
}
