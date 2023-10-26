package provider

import (
	"github.com/gomodule/redigo/redis"
	"quick_web_golang/config"
	"quick_web_golang/log"
)

type Redis struct {
	Pool *redis.Pool
}

func (r *Redis) New() *Redis {
	r.Pool = &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Get(config.RedisAddr))
			if err != nil {
				log.Fatal("Error loading redis: ", err)
				return nil, err
			}
			if _, err := c.Do("AUTH", config.Get(config.RedisPassword)); err != nil {
				_ = c.Close()

				log.Fatal("Error loading redis: ", err)
				return nil, err
			}

			return c, err
		},
	}
	return r
}

func (r *Redis) Start() {
	return
}

func (r *Redis) Close() {
	if err := r.Pool.Close(); err != nil {
		_ = log.Error(err)
	}
}
