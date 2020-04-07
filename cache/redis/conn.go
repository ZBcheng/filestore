package redis

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/gomodule/redigo/redis"
)

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
	config, err := goconfig.LoadConfigFile("db.conf")

	if err != nil {
		fmt.Println("Failed to read db.conf, err: ", err.Error())
		os.Exit(-1)
	}

	maxIdleCount, _ := config.GetValue("redis", "MaxIdle")
	maxActiveCount, _ := config.GetValue("reids", "MaxActive")
	host, _ = config.GetValue("redis", "host")

	maxIdle, _ = strconv.Atoi(maxIdleCount)
	maxActive, _ = strconv.Atoi(maxActiveCount)

	pool = newRedisPool()

}

// RedisPool : 返回redis连接池
func RedisPool() *redis.Pool {
	return pool
}
