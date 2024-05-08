package redisdb

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/nartvt/go-core/conf"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(rediConf *conf.Redis) *RedisClient {
	if len(rediConf.Username) == 0 {
		rediConf.Username = "default"
	}
	config := &redis.Options{
		Addr:     rediConf.Addr,
		Username: rediConf.Username,
		Password: rediConf.Pass,    // no password set
		DB:       int(rediConf.Db), // use default DB
		Protocol: 3,
	}
	if rediConf.Ssl {
		config.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return &RedisClient{client: redis.NewClient(config)}
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key, val string, expiredInSec int32) (string, error) {
	return r.client.Set(ctx, key, val, time.Duration(expiredInSec)*time.Second).Result()
}
