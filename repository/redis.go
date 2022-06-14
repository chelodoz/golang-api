package repository

import (
	"context"
	"errors"
	"fmt"
	"golang-api/entity"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisTokenRepository struct {
	Client *redis.Client
}

func NewRedisCache(host string, db int, exp time.Duration) entity.TokenRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       db,
	})
	return &redisTokenRepository{
		Client: client,
	}
}

// DeleteRefreshToken implements entity.TokenRepository
func (redisRepository *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := redisRepository.Client.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return err
	}

	// Val returns count of deleted keys.
	// If no key was deleted, the refresh token is invalid
	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userID/tokenID: %s/%s does not exist\n", userID, tokenID)
		return errors.New("invalid refresh token")
	}

	return nil
}

func (redisRepository *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s*", userID)

	iter := redisRepository.Client.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := redisRepository.Client.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("Failed to delete refresh token: %s\n", iter.Val())
			failCount++
		}
	}

	// check last value
	if err := iter.Err(); err != nil {
		log.Printf("Failed to delete refresh token: %s\n", iter.Val())
	}

	if failCount > 0 {
		return errors.New("invalid refresh token")
	}

	return nil
}

func (redisRepository *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Time) error {
	now := time.Now()
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := redisRepository.Client.Set(ctx, key, 0, expiresIn.Sub(now)).Err(); err != nil {
		return errors.New("could not SET refresh token to redis for userID/tokenID")
	}
	return nil
}

// func (cache *redisTokenRepository) Get(key string) *entity.TokenDetails {
// 	var tokenDetails entity.TokenDetails

// 	client := cache.getClient()
// 	val, err := client.Get(key).Result()
// 	if err != nil {
// 		return nil
// 	}

// 	err = json.Unmarshal([]byte(val), &tokenDetails)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return &tokenDetails
// }

// func (cache *redisTokenRepository) Set(tokenDetails *entity.TokenDetails) {
// 	client := cache.getClient()

// 	now := time.Now()

// 	json, err := json.Marshal(tokenDetails)

// 	if err != nil {
// 		panic(err)
// 	}

// 	err = client.Set(tokenDetails.AccessToken, json, tokenDetails.AccessTokenExpiresAt.Sub(now)).Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = client.Set(tokenDetails.RefreshToken, json, tokenDetails.RefreshTokenExpiresAt.Sub(now)).Err()
// 	if err != nil {
// 		panic(err)
// 	}
// }
