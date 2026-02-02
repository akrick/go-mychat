package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"akrick.com/mychat/admin/backend/cache"
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境需要验证
	},
}

// Client WebSocket客户端
type Client struct {
	ID       uint
	Conn     *websocket.Conn
	Send     chan []byte
	SessionID *uint
}

// Hub WebSocket连接中心
type Hub struct {
	clients    map[uint]*Client // userID -> Client
	sessions   map[uint]map[uint]*Client // sessionID -> map[userID]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

var globalHub *Hub

// InitHub 初始化Hub
func InitHub() {
	globalHub = &Hub{
		clients:    make(map[uint]*Client),
		sessions:   make(map[uint]map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
	go globalHub.Run()
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("客户端注册: userID=%d", client.ID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
				log.Printf("客户端断开: userID=%d", client.ID)
			}
			
			// 从会话中移除
			if client.SessionID != nil {
				if sessionClients, ok := h.sessions[*client.SessionID]; ok {
					delete(sessionClients, client.ID)
					if len(sessionClients) == 0 {
						delete(h.sessions, *client.SessionID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// HandleWebSocket 处理WebSocket连接
func HandleWebSocket(c *gin.Context) {
	if globalHub == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket Hub未初始化"})
		return
	}

	// 从URL参数获取token
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少token"})
		return
	}

	// 验证token
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token无效"})
		return
	}

	// 升级HTTP连接到WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	client := &Client{
		ID:   claims.UserID,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	// 注册客户端
	globalHub.register <- client

	// 启动读写协程
	go client.readPump()
	go client.writePump()
}

// WSMessage WebSocket消息结构
type WSMessage struct {
	Type      string                 `json:"type"`      // 消息类型: join, message, leave, billing, end
	SessionID uint                   `json:"session_id"` // 会话ID
	Data      map[string]interface{} `json:"data"`      // 消息数据
}

// 读取消息
func (c *Client) readPump() {
	defer func() {
		globalHub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket读取错误: %v", err)
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("消息解析失败: %v", err)
			continue
		}

		// 处理消息
		switch wsMsg.Type {
		case "join":
			c.handleJoin(wsMsg)
		case "message":
			c.handleMessage(wsMsg)
		case "leave":
			c.handleLeave(wsMsg)
		case "ping":
			c.handlePing()
		case "typing":
			messageHandler.handleTyping(c, wsMsg)
		case "typing_stop":
			messageHandler.handleTypingStop(c, wsMsg)
		case "read":
			messageHandler.handleRead(c, wsMsg)
		}
	}
}

// 写入消息
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 处理加入会话
func (c *Client) handleJoin(wsMsg WSMessage) {
	sessionID := wsMsg.SessionID

	// 查询会话
	var session models.ChatSession
	ctx := context.Background()

	// 先从缓存获取
	if sessionData, err := cache.GetChatSessionFromCache(ctx, sessionID); err == nil {
		session = *sessionData
	} else {
		// 缓存未命中，从数据库查询
		if err := database.DB.Where("id = ?", sessionID).First(&session).Error; err != nil {
			c.sendError("会话不存在")
			return
		}

		// 缓存会话数据
		cache.SetChatSessionToCache(ctx, &session, 5*time.Minute)
	}

	// 检查权限
	if session.UserID != c.ID && session.CounselorID != c.ID {
		c.sendError("无权加入此会话")
		return
	}

	// 检查会话状态
	if session.Status != 0 && session.Status != 1 {
		c.sendError("会话已结束")
		return
	}

	// 设置会话ID
	c.SessionID = &sessionID

	// 加入会话
	globalHub.mu.Lock()
	if _, ok := globalHub.sessions[sessionID]; !ok {
		globalHub.sessions[sessionID] = make(map[uint]*Client)
	}
	globalHub.sessions[sessionID][c.ID] = c
	globalHub.mu.Unlock()

	// 如果会话状态是待开始，且双方都已加入，则开始会话
	if session.Status == 0 {
		if len(globalHub.sessions[sessionID]) >= 2 {
			c.startSession(sessionID)
		}
	}

	// 更新会话管理器
	if session.Status == 1 {
		var counselor models.Counselor
		database.DB.First(&counselor, session.CounselorID)
		sessionManager.StartSession(sessionID, session.UserID, session.CounselorID, counselor.Price)
	}

	// 发送加入成功消息
	c.sendMessage("join_success", gin.H{
		"session_id": sessionID,
		"status":     session.Status,
	})
}

// 开始会话
func (c *Client) startSession(sessionID uint) {
	now := time.Now()

	// 更新会话状态
	database.DB.Model(&models.ChatSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"status":     1,
		"start_time": now,
	})

	// 查询会话和订单
	var session models.ChatSession
	database.DB.Preload("Order").Preload("Counselor").First(&session, sessionID)

	// 清除缓存
	ctx := context.Background()
	cache.DeleteChatSessionCache(ctx, sessionID)

	// 通知会话内所有客户端
	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	if clients, ok := globalHub.sessions[sessionID]; ok {
		msg, _ := json.Marshal(WSMessage{
			Type:      "session_start",
			SessionID: sessionID,
			Data: gin.H{
				"start_time": now,
				"price":      session.Counselor.Price,
			},
		})

		for _, client := range clients {
			select {
			case client.Send <- msg:
			default:
			}
		}
	}
}

// 处理聊天消息
func (c *Client) handleMessage(wsMsg WSMessage) {
	sessionID := wsMsg.SessionID
	content, _ := wsMsg.Data["content"].(string)
	contentType, _ := wsMsg.Data["content_type"].(string)
	if contentType == "" {
		contentType = "text"
	}
	fileURL, _ := wsMsg.Data["file_url"].(string)

	// 查询会话
	var session models.ChatSession
	if err := database.DB.Where("id = ?", sessionID).First(&session).Error; err != nil {
		c.sendError("会话不存在")
		return
	}

	// 检查权限
	if session.UserID != c.ID && session.CounselorID != c.ID {
		c.sendError("无权发送消息")
		return
	}

	// 确定发送者类型
	senderType := "user"
	if session.CounselorID == c.ID {
		senderType = "counselor"
	}

	// 创建消息记录
	message := models.ChatMessage{
		SessionID:   sessionID,
		SenderID:    c.ID,
		SenderType:  senderType,
		ContentType: contentType,
		Content:     content,
		FileURL:     fileURL,
		IsRead:      false,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		c.sendError("消息保存失败")
		return
	}

	// 广播消息给会话内其他客户端
	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	if clients, ok := globalHub.sessions[sessionID]; ok {
		msg, _ := json.Marshal(WSMessage{
			Type:      "message",
			SessionID: sessionID,
			Data: gin.H{
				"message_id":   message.ID,
				"sender_id":    message.SenderID,
				"sender_type":  message.SenderType,
				"content_type": message.ContentType,
				"content":      message.Content,
				"file_url":     message.FileURL,
				"created_at":   message.CreatedAt,
			},
		})

		for userID, client := range clients {
			if userID != c.ID {
				select {
				case client.Send <- msg:
				default:
				}
			}
		}
	}
}

// 处理离开会话
func (c *Client) handleLeave(wsMsg WSMessage) {
	sessionID := wsMsg.SessionID

	// 查询会话
	var session models.ChatSession
	if err := database.DB.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return
	}

	// 只有咨询师可以结束会话
	if session.CounselorID != c.ID {
		c.sendError("只有咨询师可以结束会话")
		return
	}

	// 检查会话状态
	if session.Status != 1 {
		c.sendError("会话未进行中")
		return
	}

	// 结束会话并计费
	c.endSession(sessionID, session)

	// 更新会话管理器
	sessionManager.EndSession(sessionID)

	// 通知会话内所有客户端
	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	if clients, ok := globalHub.sessions[sessionID]; ok {
		msg, _ := json.Marshal(WSMessage{
			Type:      "session_end",
			SessionID: sessionID,
			Data: gin.H{
				"ended_by": c.ID,
			},
		})

		for _, client := range clients {
			select {
			case client.Send <- msg:
			default:
			}
		}
	}
}

// 结束会话并计费
func (c *Client) endSession(sessionID uint, session models.ChatSession) {
	now := time.Now()
	
	// 计算时长
	duration := int(now.Sub(*session.StartTime).Seconds())
	
	// 查询咨询师获取单价
	var counselor models.Counselor
	database.DB.First(&counselor, session.CounselorID)
	
	pricePerMinute := counselor.Price
	
	// 计算总金额（按分钟向上取整）
	durationMinutes := (duration + 59) / 60
	if durationMinutes < 1 {
		durationMinutes = 1
	}
	totalAmount := float64(durationMinutes) * pricePerMinute
	
	// 计算平台费用和咨询师费用
	platformFee := totalAmount * 0.30
	counselorFee := totalAmount * 0.70
	
	// 更新会话
	database.DB.Model(&session).Updates(map[string]interface{}{
		"status":      2,
		"end_time":    now,
		"duration":    duration,
		"price":       pricePerMinute,
		"total_amount": totalAmount,
	})
	
	// 创建计费记录
	billing := models.ChatBilling{
		SessionID:      sessionID,
		OrderID:        session.OrderID,
		UserID:         session.UserID,
		CounselorID:    session.CounselorID,
		Duration:       duration,
		PricePerMinute: pricePerMinute,
		TotalAmount:    totalAmount,
		PlatformFee:    platformFee,
		CounselorFee:   counselorFee,
		Status:         0, // 待结算
	}
	database.DB.Create(&billing)
	
	// 更新咨询师账户
	var account models.CounselorAccount
	err := database.DB.Where("counselor_id = ?", session.CounselorID).First(&account).Error
	if err != nil {
		account = models.CounselorAccount{
			CounselorID: session.CounselorID,
			Balance:     0,
		}
		database.DB.Create(&account)
	}

	database.DB.Model(&account).Updates(map[string]interface{}{
		"total_income": account.TotalIncome + counselorFee,
		"balance":      account.Balance + counselorFee,
	})

	// 清除缓存
	ctx := context.Background()
	cache.DeleteChatSessionCache(ctx, sessionID)
	cache.DeleteCounselorAccountCache(ctx, session.CounselorID)
	
	// 发送计费信息给用户
	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()
	
	if clients, ok := globalHub.sessions[sessionID]; ok {
		billingMsg, _ := json.Marshal(WSMessage{
			Type:      "billing",
			SessionID: sessionID,
			Data: gin.H{
				"duration":       duration,
				"duration_minutes": durationMinutes,
				"price_per_minute": pricePerMinute,
				"total_amount":    totalAmount,
				"platform_fee":    platformFee,
				"counselor_fee":   counselorFee,
			},
		})
		
		for userID, client := range clients {
			if session.UserID == userID {
				select {
				case client.Send <- billingMsg:
				default:
				}
			}
		}
	}
}

// 处理ping
func (c *Client) handlePing() {
	c.sendMessage("pong", gin.H{"timestamp": time.Now().Unix()})

	// 更新会话最后ping时间
	if c.SessionID != nil {
		sessionManager.UpdateLastPing(*c.SessionID)
	}
}

// 发送消息
func (c *Client) sendMessage(msgType string, data map[string]interface{}) {
	msg, _ := json.Marshal(WSMessage{
		Type: msgType,
		Data: data,
	})
	select {
	case c.Send <- msg:
	default:
	}
}

// 发送错误消息
func (c *Client) sendError(errorMsg string) {
	c.sendMessage("error", gin.H{"error": errorMsg})
}

// BroadcastToSession 向指定会话广播消息
func BroadcastToSession(sessionID uint, message []byte) {
	if globalHub == nil {
		return
	}

	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	if clients, ok := globalHub.sessions[sessionID]; ok {
		for _, client := range clients {
			select {
			case client.Send <- message:
			default:
			}
		}
	}
}

// BroadcastToAll 向所有在线用户广播消息
func BroadcastToAll(message []byte) {
	if globalHub == nil {
		return
	}
	globalHub.broadcast <- message
}

// GetOnlineUsers 获取在线用户列表
func GetOnlineUsers() []uint {
	if globalHub == nil {
		return []uint{}
	}

	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	users := make([]uint, 0, len(globalHub.clients))
	for userID := range globalHub.clients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline 检查用户是否在线
func IsUserOnline(userID uint) bool {
	if globalHub == nil {
		return false
	}

	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	_, ok := globalHub.clients[userID]
	return ok
}

// GetSessionParticipants 获取会话参与者
func GetSessionParticipants(sessionID uint) []uint {
	if globalHub == nil {
		return nil
	}

	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	if clients, ok := globalHub.sessions[sessionID]; ok {
		participants := make([]uint, 0, len(clients))
		for userID := range clients {
			participants = append(participants, userID)
		}
		return participants
	}
	return nil
}

// BuildRevokeMessage 构建撤回消息
func BuildRevokeMessage(messageID string) []byte {
	msg, _ := json.Marshal(WSMessage{
		Type: "message_revoked",
		Data: gin.H{
			"message_id": messageID,
		},
	})
	return msg
}

// SendToUser 发送消息给指定用户
func SendToUser(userID uint, message []byte) bool {
	if globalHub == nil {
		return false
	}

	globalHub.mu.RLock()
	client, ok := globalHub.clients[userID]
	globalHub.mu.RUnlock()

	if !ok {
		return false
	}

	select {
	case client.Send <- message:
		return true
	default:
		return false
	}
}
