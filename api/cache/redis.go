package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var (
	Rdb     *redis.Client
	Request singleflight.Group
)

type UserInfoCache struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
}

func InitRedis() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis连接失败: %w", err)
	}

	return nil
}

func CloseRedis() error {
	return Rdb.Close()
}

func getUserCacheKey(userID uint) string {
	return fmt.Sprintf("user:info:%d", userID)
}

// 获取用户信息（带缓存和singleflight防穿透）
func GetUserInfoWithCache(ctx context.Context, userID uint) (*UserInfoCache, error) {
	cacheKey := getUserCacheKey(userID)

	// 使用singleflight防止缓存穿透
	result, err, _ := Request.Do(cacheKey, func() (any, error) {
		// 1. 先从缓存获取
		cachedData, err := Rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			var userInfo UserInfoCache
			if err := json.Unmarshal([]byte(cachedData), &userInfo); err == nil {
				return &userInfo, nil
			}
		}

		// 2. 缓存未命中，从数据库查询
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			// 查询失败，缓存空值防止穿透（5分钟）
			if err := Rdb.Set(ctx, cacheKey, "", 5*time.Minute).Err(); err != nil {
				return nil, err
			}
			return nil, err
		}

		// 3. 将数据写入缓存（30分钟）
		userInfo := UserInfoCache{
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			Status:   user.Status,
		}

		data, _ := json.Marshal(userInfo)
		// 添加10分钟以内的随机过期时间，防止缓存雪崩
		ttl := 30*time.Minute + time.Duration(rand.Intn(10))*time.Minute
		if err := Rdb.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
			return nil, err
		}

		return &userInfo, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*UserInfoCache), nil
}

// 删除用户缓存
func DeleteUserCache(ctx context.Context, userID uint) error {
	return Rdb.Del(ctx, getUserCacheKey(userID)).Err()
}

// 刷新用户缓存
func RefreshUserCache(ctx context.Context, userID uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	userInfo := UserInfoCache{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Status:   user.Status,
	}

	data, _ := json.Marshal(userInfo)
	cacheKey := getUserCacheKey(userID)
	// 添加10分钟以内的随机过期时间，防止缓存雪崩
	ttl := 30*time.Minute + time.Duration(rand.Intn(10))*time.Minute
	return Rdb.Set(ctx, cacheKey, data, ttl).Err()
}
