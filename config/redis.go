package config

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	redisAddress  = ":6379"
	redisDatabase = 2
)

var (
	pool *redis.Pool
)

// XXX: make this thread-safe
// possibly using sync.Once?
func Pool() *redis.Pool {
	if pool != nil {
		return pool
	}
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisAddress)
			if err != nil {
				return nil, err
			}
			_, err = conn.Do("SELECT", redisDatabase)
			if err != nil {
				defer conn.Close()
				return nil, err
			}

			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
	return pool
}
