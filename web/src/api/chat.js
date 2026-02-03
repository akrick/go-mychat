import request from '@/utils/request'

export function startChatSession(orderId) {
  return request({
    url: `/chat/start/${orderId}`,
    method: 'post'
  })
}

/**
 * 获取订单的会话ID
 * @param {number} orderId - 订单ID
 */
export function getOrderSessionId(orderId) {
  return request({
    url: `/chat/order/${orderId}/session`,
    method: 'get'
  })
}

export function sendMessage(sessionId, data) {
  return request({
    url: `/chat/session/${sessionId}/message`,
    method: 'post',
    data
  })
}

export function getMessages(sessionId, params) {
  return request({
    url: `/chat/messages/${sessionId}`,
    method: 'get',
    params
  })
}

export function endChatSession(sessionId) {
  return request({
    url: `/chat/end/${sessionId}`,
    method: 'post'
  })
}

export function getChatSessions(params) {
  return request({
    url: '/chat/sessions',
    method: 'get',
    params
  })
}
