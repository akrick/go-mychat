import request from '@/utils/request'

// 获取聊天会话列表
export function getChatSessions(params) {
  return request({
    url: '/api/admin/chat/sessions',
    method: 'get',
    params
  })
}

// 获取聊天消息列表
export function getChatMessages(sessionId, params) {
  return request({
    url: `/api/admin/chat/sessions/${sessionId}/messages`,
    method: 'get',
    params
  })
}

// 获取聊天统计
export function getChatStatistics() {
  return request({
    url: '/api/admin/chat/statistics',
    method: 'get'
  })
}

// 搜索聊天消息
export function searchChatMessages(params) {
  return request({
    url: '/api/admin/chat/messages/search',
    method: 'get',
    params
  })
}

// 删除聊天会话
export function deleteChatSession(id) {
  return request({
    url: `/api/admin/chat/sessions/${id}`,
    method: 'delete'
  })
}
