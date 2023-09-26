package cache

import (
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// Redis contains the client/connection used by the cache manager
type Redis struct {
	Source *redis.Client
}

// NewRedis gets new instance of the redis client
func NewRedis(host string, port string, pass string, db int, useTLS bool) *Redis {
	redisOpts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
		DB:       db,
	}
	if useTLS {
		redisOpts.TLSConfig = &tls.Config{}
	}

	redisClient := redis.NewClient(redisOpts)

	return &Redis{
		Source: redisClient,
	}
}

// GetRedis gets the redis client/connection
func (r *Redis) GetRedis() *redis.Client {
	return r.Source
}

// CloseRedis closes the redis client/connection
func (r *Redis) CloseRedis() {
	err := r.Source.Close()
	if err != nil {
		return
	}
}
