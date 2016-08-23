package nuvi

import (
	"io"
	"io/ioutil"

	"gopkg.in/redis.v4"
)

//go:generate counterfeiter . RedisClient

// RedisClient connects to Redis
type RedisClient interface {
	LPush(key string, values ...interface{}) *redis.IntCmd
}

// RedisCacher caches content into redis
type RedisCacher struct {
	Key    string
	Client RedisClient
}

// Cache caches the content of io.Reader
func (cacher *RedisCacher) Cache(reader io.Reader) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	cacher.Client.LPush(cacher.Key, string(data))
}
