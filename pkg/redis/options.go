package redis

import (
	"strings"
	"time"

	"github.com/idoyudha/eshop-cart/config"
	"github.com/redis/go-redis/v9"
)

func RedisFailoverOptions(cfg config.Redis) *redis.FailoverOptions {
	sentinelAddrs := strings.Split(cfg.RedisSentinelAddrs, ",")
	return &redis.FailoverOptions{
		MasterName:    cfg.RedisMaster,
		SentinelAddrs: sentinelAddrs,
		Password:      cfg.RedisPassword,
		DB:            0,
		ReadTimeout:   time.Second * 3,
		WriteTimeout:  time.Second * 3,
	}
}
