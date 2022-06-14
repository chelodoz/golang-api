package cache

import (
	"encoding/json"
	"golang-api/entity"
	"time"

	"github.com/go-redis/redis"
)

type Cache interface {
	Set(value *entity.TokenDetails)
	Get(key string) *entity.TokenDetails
}

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) Cache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}
func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}
func (cache *redisCache) Get(key string) *entity.TokenDetails {
	var tokenDetails entity.TokenDetails

	client := cache.getClient()
	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	err = json.Unmarshal([]byte(val), &tokenDetails)
	if err != nil {
		panic(err)
	}

	return &tokenDetails
}

func (cache *redisCache) Set(tokenDetails *entity.TokenDetails) {
	client := cache.getClient()

	now := time.Now()

	json, err := json.Marshal(tokenDetails)

	if err != nil {
		panic(err)
	}

	err = client.Set(tokenDetails.AccessToken, json, tokenDetails.AccessTokenExpiresAt.Sub(now)).Err()
	if err != nil {
		panic(err)
	}

	err = client.Set(tokenDetails.RefreshToken, json, tokenDetails.RefreshTokenExpiresAt.Sub(now)).Err()
	if err != nil {
		panic(err)
	}
}
