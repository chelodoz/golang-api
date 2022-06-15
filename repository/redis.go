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

func NewRedisCache(host string, port string, db int) entity.TokenRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       db,
	})
	return &redisTokenRepository{
		Client: client,
	}
}

func (redisRepository *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userEmail string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userEmail, tokenID)

	result := redisRepository.Client.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userEmail/tokenID: %s/%s: %v\n", userEmail, tokenID, err)
		return err
	}

	// Val returns count of deleted keys.
	// If no key was deleted, the refresh token is invalid
	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userEmail/tokenID: %s/%s does not exist\n", userEmail, tokenID)
		return errors.New("invalid refresh token")
	}

	return nil
}

func (redisRepository *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userEmail string) error {
	pattern := fmt.Sprintf("%s*", userEmail)

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

func (redisRepository *redisTokenRepository) SetRefreshToken(ctx context.Context, userEmail string, tokenID string, expiresIn time.Time) error {
	now := time.Now()
	key := fmt.Sprintf("%s:%s", userEmail, tokenID)
	if err := redisRepository.Client.Set(ctx, key, 0, expiresIn.Sub(now)).Err(); err != nil {
		return errors.New("could not SET refresh token to redis for userEmail/tokenID")
	}
	return nil
}
