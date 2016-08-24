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
	LIndex(key string, index int64) *redis.StringCmd
	LLen(key string) *redis.IntCmd
}

// RedisCacher caches content into redis
type RedisCacher struct {
	Key    string
	Client RedisClient
	Logger Logger
}

// Cache caches the content of io.Reader
func (cacher *RedisCacher) Cache(reader io.Reader) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	text := string(data)
	cacher.Logger.Println("Caching data")
	// Using Redis list is not the best data structur
	// in order to keep uniqueness of the elements
	// because we need to iterate over the list
	// O(n) complexity
	length, err := cacher.Client.LLen(cacher.Key).Result()
	if err != nil {
		return
	}

	for index := int64(0); index < length; index++ {
		xmlDoc, err := cacher.Client.LIndex(cacher.Key, index).Result()
		if err != nil {
			return
		}
		if xmlDoc == text {
			return
		}
	}

	cacher.Client.LPush(cacher.Key, text)
}
