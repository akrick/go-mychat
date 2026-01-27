package websocket

import (
	"encoding/json"
	"log"
	"time"

	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"

	"github.com/gin-gonic/gin"
)

// MessageHandler 消息处理器
type MessageHandler struct{}

var messageHandler = &MessageHandler{}

// HandleMessage 处理WebSocket消息（未使用，方法已在hub.go中直接实现）
func HandleMessage(c *Client, wsMsg WSMessage) {
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
	case "read":
		messageHandler.handleRead(c, wsMsg)
	case "typing_stop":
		messageHandler.handleTypingStop(c, wsMsg)
	default:
		log.Printf("未知消息类型: %s", wsMsg.Type)
		c.sendError("未知消息类型")
	}
}

// handleTyping 处理正在输入状态
func (mh *MessageHandler) handleTyping(c *Client, wsMsg WSMessage) {
	sessionID := wsMsg.SessionID

	// 检查权限
	var session models.ChatSession
	if err := database.DB.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return
	}

	if session.UserID != c.ID && session.CounselorID != c.ID {
		return
	}

	// 确定对方ID
	opponentID := session.UserID
	if session.UserID == c.ID {
		opponentID = session.CounselorID
	}

	// 发送正在输入通知
	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	if clients, ok := globalHub.sessions[sessionID]; ok {
		msg, _ := json.Marshal(WSMessage{
			Type:      "typing",
			SessionID: sessionID,
			Data: gin.H{
				"user_id": c.ID,
			},
		})

		if client, ok := clients[opponentID]; ok {
			select {
			case client.Send <- msg:
			default:
			}
		}
	}
}

// handleTypingStop 处理停止输入状态
func (mh *MessageHandler) handleTypingStop(c *Client, wsMsg WSMessage) {
	sessionID := wsMsg.SessionID

	// 检查权限
	var session models.ChatSession
	if err := database.DB.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return
	}

	if session.UserID != c.ID && session.CounselorID != c.ID {
		return
	}

	// 确定对方ID
	opponentID := session.UserID
	if session.UserID == c.ID {
		opponentID = session.CounselorID
	}

	// 发送停止输入通知
	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	if clients, ok := globalHub.sessions[sessionID]; ok {
		msg, _ := json.Marshal(WSMessage{
			Type:      "typing_stop",
			SessionID: sessionID,
			Data: gin.H{
				"user_id": c.ID,
			},
		})

		if client, ok := clients[opponentID]; ok {
			select {
			case client.Send <- msg:
			default:
			}
		}
	}
}

// handleRead 处理消息已读
func (mh *MessageHandler) handleRead(c *Client, wsMsg WSMessage) {
	sessionID := wsMsg.SessionID
	messageID, _ := wsMsg.Data["message_id"].(float64)

	// 查询消息
	var message models.ChatMessage
	if err := database.DB.Where("id = ? AND session_id = ?", messageID, sessionID).First(&message).Error; err != nil {
		return
	}

	// 只有接收者可以标记为已读
	var session models.ChatSession
	if err := database.DB.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return
	}

	receiverID := session.UserID
	if session.UserID == c.ID {
		receiverID = session.CounselorID
	}

	if receiverID != c.ID {
		return
	}

	// 更新已读状态
	nowTime := time.Now()
	message.IsRead = true
	message.ReadTime = &nowTime
	database.DB.Save(&message)

	// 通知发送者消息已读
	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	if clients, ok := globalHub.sessions[sessionID]; ok {
		msg, _ := json.Marshal(WSMessage{
			Type:      "message_read",
			SessionID: sessionID,
			Data: gin.H{
				"message_id": messageID,
				"read_by":    c.ID,
				"read_at":    nowTime,
			},
		})

		if client, ok := clients[message.SenderID]; ok {
			select {
			case client.Send <- msg:
			default:
			}
		}
	}
}

// BroadcastUserMessage 向指定用户广播消息
func BroadcastUserMessage(userID uint, msgType string, data map[string]interface{}) {
	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()

	if client, ok := globalHub.clients[userID]; ok {
		msg, _ := json.Marshal(WSMessage{
			Type: msgType,
			Data: data,
		})

		select {
		case client.Send <- msg:
		default:
		}
	}
}
