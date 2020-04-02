package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
)

// newRedisPool : 创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 1. 打开连接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			return c, nil
		},
	}
}

func init() {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}
