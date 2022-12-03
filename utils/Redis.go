package utils

import (
	"fmt"
	"okc/config"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func init() {
	m := config.Config().REDIS.(map[interface{}]interface{})
	// 设置redis 连接池
	pool = &redis.Pool{
		MaxIdle:     10,  // 最大空闲连接数量
		MaxActive:   0,   // 最大连接数量 0代表不限制
		IdleTimeout: 100, // 最大空闲时间,超过该秒自动回收
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", m["DATABASE_HOST"], m["DATABASE_PORT"]))
		},
	}
}

func RedisConnect() (connect redis.Conn) {
	connect = pool.Get()
	return
}
