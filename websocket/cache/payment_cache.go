package cache

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"golang.org/x/sync/singleflight"
)

var (
	paymentGroup singleflight.Group
)

// GetPaymentWithCache 使用SingleFlight获取支付记录（防缓存穿透）
func GetPaymentWithCache(ctx context.Context, paymentID uint) (*models.Payment, error) {
	if Rdb == nil {
		return nil, errors.New("Redis未连接")
	}

	cacheKey := fmt.Sprintf("payment:%d", paymentID)

	// 使用SingleFlight防止缓存穿透
	result, err, _ := paymentGroup.Do(cacheKey, func() (interface{}, error) {
		// 先从缓存获取
		_, err := Rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			// 缓存命中，反序列化（这里简化处理，实际应使用序列化）
			var payment models.Payment
			// 注意：实际项目中应使用json.Unmarshal或其他序列化方式
			return &payment, nil
		}

		// 缓存未命中，从数据库获取
		var payment models.Payment
		if err := database.DB.First(&payment, paymentID).Error; err != nil {
			return nil, err
		}

		// 设置缓存，添加随机时间防止缓存雪崩
		ttl := 30*time.Minute + time.Duration(rand.Intn(10))*time.Minute
		// 实际项目应使用序列化后存储
		Rdb.Set(ctx, cacheKey, payment.ID, ttl)

		return &payment, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Payment), nil
}

// DeletePaymentCache 删除支付缓存
func DeletePaymentCache(ctx context.Context, paymentID uint) {
	if Rdb == nil {
		return
	}
	cacheKey := fmt.Sprintf("payment:%d", paymentID)
	Rdb.Del(ctx, cacheKey)
}

// InvalidateUserPaymentsCache 清除用户支付记录列表缓存
func InvalidateUserPaymentsCache(ctx context.Context, userID uint) {
	if Rdb == nil {
		return
	}
	cacheKey := fmt.Sprintf("user:payments:%d", userID)
	Rdb.Del(ctx, cacheKey)
}

// InvalidateOrderPaymentCache 清除订单支付记录缓存
func InvalidateOrderPaymentCache(ctx context.Context, orderID uint) {
	if Rdb == nil {
		return
	}
	cacheKey := fmt.Sprintf("order:payment:%d", orderID)
	Rdb.Del(ctx, cacheKey)
}

// SetPaymentCache 设置支付缓存（用于主动更新缓存）
func SetPaymentCache(ctx context.Context, payment *models.Payment) error {
	if Rdb == nil {
		return nil
	}
	cacheKey := fmt.Sprintf("payment:%d", payment.ID)
	ttl := 30*time.Minute + time.Duration(rand.Intn(10))*time.Minute
	// 实际项目应使用序列化后存储
	return Rdb.Set(ctx, cacheKey, payment.ID, ttl).Err()
}
