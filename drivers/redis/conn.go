package drivers

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/zbcheng/filestore/conf"
)

var config = conf.GetConfig()

var (
	pool      *redis.Pool
	maxIdle   int
	maxActive int
	host      string
)

// newRedisPool : 创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 1. 打开连接
			c, err := redis.Dial("tcp", host)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			return c, nil
		},
	}
}

func init() {
	maxIdle = config.RdConf.MaxIdle
	maxActive = config.RdConf.MaxActive

	host = config.RdConf.Host
	pool = newRedisPool()

}

// RedisPool : 返回redis连接池
func RedisPool() *redis.Pool {
	return pool
}
