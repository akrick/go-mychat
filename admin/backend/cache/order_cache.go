package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"golang.org/x/sync/singleflight"
)

type OrderCacheData struct {
	ID           uint       `json:"id"`
	OrderNo      string     `json:"order_no"`
	UserID       uint       `json:"user_id"`
	CounselorID  uint       `json:"counselor_id"`
	Duration     int        `json:"duration"`
	Amount       float64    `json:"amount"`
	Status       int        `json:"status"`
	ScheduleTime time.Time  `json:"schedule_time"`
	Notes        string     `json:"notes"`
	PayTime      *time.Time `json:"pay_time"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CounselorCacheData struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Specialty string    `json:"specialty"`
	Price     float64   `json:"price"`
	YearsExp  int       `json:"years_exp"`
	Rating    float64   `json:"rating"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 订单缓存
var (
	OrderGroup      singleflight.Group
	CounselorGroup  singleflight.Group
)

// GetOrderWithCache 获取订单详情（带缓存和singleflight）
func GetOrderWithCache(ctx context.Context, orderID uint) (*OrderCacheData, error) {
	cacheKey := fmt.Sprintf("order:detail:%d", orderID)

	result, err, _ := OrderGroup.Do(cacheKey, func() (interface{}, error) {
		// 1. 先从缓存获取
		cachedData, err := Rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			var order OrderCacheData
			if err := json.Unmarshal([]byte(cachedData), &order); err == nil {
				return &order, nil
			}
		}

		// 2. 缓存未命中，从数据库查询
		var order models.Order
		if err := database.DB.First(&order, orderID).Error; err != nil {
			// 查询失败，缓存空值防止穿透（5分钟）
			Rdb.Set(ctx, cacheKey, "", 5*time.Minute)
			return nil, err
		}

		// 3. 将数据写入缓存（30分钟 + 随机0-10分钟）
		orderData := OrderCacheData{
			ID:           order.ID,
			OrderNo:      order.OrderNo,
			UserID:       order.UserID,
			CounselorID:  order.CounselorID,
			Duration:     order.Duration,
			Amount:       order.Amount,
			Status:       order.Status,
			ScheduleTime: order.ScheduleTime,
			Notes:        order.Notes,
			PayTime:      order.PayTime,
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		}

		data, _ := json.Marshal(orderData)
		ttl := 30*time.Minute + time.Duration(rand.Intn(10))*time.Minute
		Rdb.Set(ctx, cacheKey, data, ttl)

		return &orderData, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*OrderCacheData), nil
}

// GetCounselorWithCache 获取咨询师详情（带缓存和singleflight）
func GetCounselorWithCache(ctx context.Context, counselorID uint) (*CounselorCacheData, error) {
	cacheKey := fmt.Sprintf("counselor:detail:%d", counselorID)

	result, err, _ := CounselorGroup.Do(cacheKey, func() (interface{}, error) {
		// 1. 先从缓存获取
		cachedData, err := Rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			var counselor CounselorCacheData
			if err := json.Unmarshal([]byte(cachedData), &counselor); err == nil {
				return &counselor, nil
			}
		}

		// 2. 缓存未命中，从数据库查询
		var counselor models.Counselor
		if err := database.DB.First(&counselor, counselorID).Error; err != nil {
			// 查询失败，缓存空值防止穿透（5分钟）
			Rdb.Set(ctx, cacheKey, "", 5*time.Minute)
			return nil, err
		}

		// 3. 将数据写入缓存（30分钟 + 随机0-10分钟）
		counselorData := CounselorCacheData{
			ID:        counselor.ID,
			Name:      counselor.Name,
			Title:     counselor.Title,
			Avatar:    counselor.Avatar,
			Bio:       counselor.Bio,
			Specialty: counselor.Specialty,
			Price:     counselor.Price,
			YearsExp:  counselor.YearsExp,
			Rating:    counselor.Rating,
			Status:    counselor.Status,
			CreatedAt: counselor.CreatedAt,
			UpdatedAt: counselor.UpdatedAt,
		}

		data, _ := json.Marshal(counselorData)
		ttl := 30*time.Minute + time.Duration(rand.Intn(10))*time.Minute
		Rdb.Set(ctx, cacheKey, data, ttl)

		return &counselorData, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*CounselorCacheData), nil
}

// DeleteOrderCache 删除订单缓存
func DeleteOrderCache(ctx context.Context, orderID uint) error {
	cacheKey := fmt.Sprintf("order:detail:%d", orderID)
	return Rdb.Del(ctx, cacheKey).Err()
}

// DeleteCounselorCache 删除咨询师缓存
func DeleteCounselorCache(ctx context.Context, counselorID uint) error {
	cacheKey := fmt.Sprintf("counselor:detail:%d", counselorID)
	return Rdb.Del(ctx, cacheKey).Err()
}

// InvalidateUserOrdersCache 删除用户订单列表缓存（当用户创建/更新订单时调用）
func InvalidateUserOrdersCache(ctx context.Context, userID uint) {
	// 删除该用户的所有订单缓存
	pattern := fmt.Sprintf("order:user:%d:*", userID)
	keys, _ := Rdb.Keys(ctx, pattern).Result()
	if len(keys) > 0 {
		Rdb.Del(ctx, keys...)
	}
}

// InvalidateCounselorOrdersCache 删除咨询师订单列表缓存
func InvalidateCounselorOrdersCache(ctx context.Context, counselorID uint) {
	// 删除该咨询师的所有订单缓存
	pattern := fmt.Sprintf("order:counselor:%d:*", counselorID)
	keys, _ := Rdb.Keys(ctx, pattern).Result()
	if len(keys) > 0 {
		Rdb.Del(ctx, keys...)
	}
}
