package repository

import (
	"context"
	"encoding/json"
	"go_project_template/internal/user/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type IUserRedisRepository interface {
	SetUserOTP(ctx context.Context, key string, userOTP model.UserOTPVerification, expiration time.Duration) error
	GetUserOTP(ctx context.Context, key string) (model.UserOTPVerification, error)
	RemoveUserOTP(ctx context.Context, key string) error
}

type UserRedisRepository struct {
	redisClient *redis.Client
}

func NewUserRedisRepository(redisClient *redis.Client) *UserRedisRepository {
	return &UserRedisRepository{
		redisClient: redisClient,
	}
}

func (rc *UserRedisRepository) SetUserOTP(ctx context.Context, key string, userOTP model.UserOTPVerification, expiration time.Duration) error {

	userOTPBytes, err := json.Marshal(userOTP)

	if err != nil {
		return err
	}

	err = rc.redisClient.Set(
		ctx,
		key,
		userOTPBytes,
		expiration,
	).Err()

	return err
}

func (rc *UserRedisRepository) GetUserOTP(ctx context.Context, key string) (model.UserOTPVerification, error) {
	var userOTP model.UserOTPVerification
	_userOTP, err := rc.redisClient.Get(ctx, key).Result()

	if err != nil {
		return userOTP, err
	}

	if err := json.Unmarshal([]byte(_userOTP), &userOTP); err != nil {
		return userOTP, err
	}

	return userOTP, nil
}

func (rc *UserRedisRepository) RemoveUserOTP(ctx context.Context, key string) error {
	err := rc.redisClient.Del(ctx, key).Err()

	return err
}
