package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

// newRedisPool : 创建reids连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err == nil {
				return nil, err
			}

			return c, nil
		},
	}
}

func init() {
	pool = newRedisPool()
}

// RedisPool : 暴露redis pool
func RedisPool() *redis.Pool {
	return pool
}
