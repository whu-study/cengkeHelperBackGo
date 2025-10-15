package database

import (
	"cengkeHelperBackGo/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Conf.Redis.Host, config.Conf.Redis.Port),
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败: %v", err)
		log.Println("注意: 邮箱验证码功能需要Redis支持，请启动Redis服务")
		RedisClient = nil
	} else {
		log.Println("Redis连接成功")
	}
}

// SetEmailCode 存储邮箱验证码，设置5分钟过期时间
func SetEmailCode(email, code string) error {
	if RedisClient == nil {
		return fmt.Errorf("redis未连接")
	}

	ctx := context.Background()
	key := fmt.Sprintf("email_code:%s", email)
	return RedisClient.Set(ctx, key, code, 5*time.Minute).Err()
}

// GetEmailCode 获取邮箱验证码
func GetEmailCode(email string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("redis未连接")
	}

	ctx := context.Background()
	key := fmt.Sprintf("email_code:%s", email)
	return RedisClient.Get(ctx, key).Result()
}

// DeleteEmailCode 删除邮箱验证码
func DeleteEmailCode(email string) error {
	if RedisClient == nil {
		return fmt.Errorf("redis未连接")
	}

	ctx := context.Background()
	key := fmt.Sprintf("email_code:%s", email)
	return RedisClient.Del(ctx, key).Err()
}
