package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	host    string
	db      int
	expires time.Duration
	client  *redis.Client
}

func NewRedisCache(host string, db int, expires time.Duration) GPTMessageCache {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return &RedisCache{
		host:    host,
		db:      db,
		expires: expires * time.Minute,
		client:  client,
	}
}

func (c *RedisCache) Set(key string, value string) {
	hash := md5.Sum([]byte(key))

	c.client.Set(hex.EncodeToString(hash[:]), value, c.expires*time.Second)
}

func (c *RedisCache) Get(key string) string {
	hash := md5.Sum([]byte(key))

	val, err := c.client.Get(hex.EncodeToString(hash[:])).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		panic(err)
	}

	return val
}
