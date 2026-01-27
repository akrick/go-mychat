package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"akrick.com/mychat/admin/backend/models"
)

// GetChatSessionFromCache 从缓存获取聊天会话
func GetChatSessionFromCache(ctx context.Context, sessionID uint) (*models.ChatSession, error) {
	cacheKey := fmt.Sprintf("chat:session:%d", sessionID)
	cachedData, err := Rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var session models.ChatSession
	if err := json.Unmarshal([]byte(cachedData), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// SetChatSessionToCache 设置聊天会话到缓存
func SetChatSessionToCache(ctx context.Context, session *models.ChatSession, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("chat:session:%d", session.ID)
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return Rdb.Set(ctx, cacheKey, sessionData, expiration).Err()
}

// GetCounselorAccountFromCache 从缓存获取咨询师账户
func GetCounselorAccountFromCache(ctx context.Context, counselorID uint) (*models.CounselorAccount, error) {
	cacheKey := fmt.Sprintf("counselor:account:%d", counselorID)
	cachedData, err := Rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var account models.CounselorAccount
	if err := json.Unmarshal([]byte(cachedData), &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// SetCounselorAccountToCache 设置咨询师账户到缓存
func SetCounselorAccountToCache(ctx context.Context, account *models.CounselorAccount, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("counselor:account:%d", account.CounselorID)
	accountData, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return Rdb.Set(ctx, cacheKey, accountData, expiration).Err()
}

// GetOnlineStatus 获取用户在线状态
func GetOnlineStatus(ctx context.Context, userID uint) bool {
	cacheKey := fmt.Sprintf("user:online:%d", userID)
	_, err := Rdb.Get(ctx, cacheKey).Result()
	return err == nil
}

// SetOnlineStatus 设置用户在线状态
func SetOnlineStatus(ctx context.Context, userID uint) {
	cacheKey := fmt.Sprintf("user:online:%d", userID)
	Rdb.Set(ctx, cacheKey, "1", 5*time.Minute)
}

// ClearOnlineStatus 清除用户在线状态
func ClearOnlineStatus(ctx context.Context, userID uint) {
	cacheKey := fmt.Sprintf("user:online:%d", userID)
	Rdb.Del(ctx, cacheKey)
}

// DeleteChatSessionCache 删除聊天会话缓存
func DeleteChatSessionCache(ctx context.Context, sessionID uint) {
	cacheKey := fmt.Sprintf("chat:session:%d", sessionID)
	Rdb.Del(ctx, cacheKey)
}

// DeleteCounselorAccountCache 删除咨询师账户缓存
func DeleteCounselorAccountCache(ctx context.Context, counselorID uint) {
	cacheKey := fmt.Sprintf("counselor:account:%d", counselorID)
	Rdb.Del(ctx, cacheKey)
}

