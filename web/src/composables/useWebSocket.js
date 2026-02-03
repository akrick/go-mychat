import { ref, onUnmounted } from 'vue'

/**
 * WebSocket 连接 Composable
 * @param {string} sessionId - 会话ID
 * @param {Object} options - 配置选项
 * @param {Function} options.onMessage - 消息回调
 * @param {Function} options.onOpen - 连接成功回调
 * @param {Function} options.onClose - 连接关闭回调
 * @param {Function} options.onError - 错误回调
 * @param {number} options.reconnectInterval - 重连间隔（毫秒）
 * @param {number} options.maxReconnectAttempts - 最大重连次数
 */
export function useWebSocket(sessionId, options = {}) {
  const {
    onMessage,
    onOpen,
    onClose,
    onError,
    reconnectInterval = 3000,
    maxReconnectAttempts = 5
  } = options

  const ws = ref(null)
  const connected = ref(false)
  const reconnectAttempts = ref(0)
  const reconnectTimer = ref(null)

  // 获取WebSocket URL
  const getWebSocketUrl = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const token = localStorage.getItem('token')
    // 开发环境使用代理
    const wsUrl = `${protocol}//${host}/ws/chat/${sessionId}?token=${token}`
    console.log('WebSocket URL:', wsUrl)
    return wsUrl
  }

  // 连接WebSocket
  const connect = () => {
    try {
      ws.value = new WebSocket(getWebSocketUrl())

      ws.value.onopen = () => {
        console.log('WebSocket连接成功')
        connected.value = true
        reconnectAttempts.value = 0
        onOpen?.()
      }

      ws.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          onMessage?.(data)
        } catch (error) {
          console.error('解析WebSocket消息失败:', error)
        }
      }

      ws.value.onerror = (error) => {
        console.error('WebSocket错误:', error)
        connected.value = false
        onError?.(error)
      }

      ws.value.onclose = () => {
        console.log('WebSocket连接关闭')
        connected.value = false
        onClose?.()

        // 自动重连
        if (reconnectAttempts.value < maxReconnectAttempts) {
          reconnectAttempts.value++
          console.log(`尝试重连 (${reconnectAttempts.value}/${maxReconnectAttempts})`)
          reconnectTimer.value = setTimeout(() => {
            connect()
          }, reconnectInterval)
        }
      }
    } catch (error) {
      console.error('创建WebSocket连接失败:', error)
      onError?.(error)
    }
  }

  // 发送消息
  const send = (data) => {
    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
      console.error('WebSocket未连接')
      return false
    }

    try {
      ws.value.send(JSON.stringify(data))
      return true
    } catch (error) {
      console.error('发送WebSocket消息失败:', error)
      return false
    }
  }

  // 关闭连接
  const close = () => {
    if (reconnectTimer.value) {
      clearTimeout(reconnectTimer.value)
      reconnectTimer.value = null
    }

    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    connected.value = false
  }

  // 组件卸载时关闭连接
  onUnmounted(() => {
    close()
  })

  return {
    ws,
    connected,
    connect,
    send,
    close
  }
}

/**
 * 聊天室 WebSocket Composable
 * @param {string} sessionId - 会话ID
 */
export function useChatWebSocket(sessionId) {
  const messages = ref([])
  const { connected, connect, send, close } = useWebSocket(sessionId, {
    onMessage: (data) => {
      if (data.type === 'message') {
        messages.value.push(data.payload)
      } else if (data.type === 'typing') {
        // 处理输入状态
      } else if (data.type === 'session_end') {
        // 处理会话结束
      }
    },
    onOpen: () => {
      console.log('聊天室连接成功')
    },
    onClose: () => {
      console.log('聊天室连接关闭')
    },
    onError: (error) => {
      console.error('聊天室连接错误:', error)
    }
  })

  const sendTextMessage = (content) => {
    return send({
      type: 'message',
      payload: {
        content_type: 'text',
        content: content
      }
    })
  }

  const sendTypingStatus = (isTyping) => {
    return send({
      type: 'typing',
      payload: {
        is_typing: isTyping
      }
    })
  }

  return {
    connected,
    messages,
    connect,
    sendTextMessage,
    sendTypingStatus,
    close
  }
}
