package cache

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, expires time.Duration) GPTMessageCache {
	return &RedisCache{
		host:    host,
		db:      db,
		expires: expires * time.Minute,
	}
}

func (c *RedisCache) GetClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.host,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       c.db,
	})
}

func (c *RedisCache) Set(key string, value string) {
	client := c.GetClient()
	hash := md5.Sum([]byte(key))

	client.Set(hex.EncodeToString(hash[:]), value, c.expires*time.Second)
}

func (c *RedisCache) Get(key string) string {
	client := c.GetClient()
	hash := md5.Sum([]byte(key))

	val, err := client.Get(hex.EncodeToString(hash[:])).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		panic(err)
	}

	return val
}
