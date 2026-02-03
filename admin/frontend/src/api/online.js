import request from '@/utils/request'

/**
 * 获取在线用户列表
 * @param {Object} params - 查询参数
 */
export function getOnlineUsers(params) {
  return request({
    url: '/admin/online/users/detailed',
    method: 'get',
    params
  })
}

/**
 * 踢用户下线
 * @param {number} userId - 用户ID
 */
export function kickUser(userId) {
  return request({
    url: `/admin/online/users/${userId}/kick`,
    method: 'post'
  })
}

/**
 * 禁言/解禁用户
 * @param {Object} data - 禁言信息
 * @param {number} data.user_id - 用户ID
 * @param {boolean} data.is_muted - 是否禁言
 */
export function muteUser(data) {
  return request({
    url: '/admin/online/mute',
    method: 'post',
    data
  })
}

/**
 * 广播系统消息
 * @param {Object} data - 消息信息
 * @param {string} data.content - 消息内容
 */
export function broadcastMessage(data) {
  return request({
    url: '/admin/broadcast',
    method: 'post',
    data
  })
}

/**
 * 获取在线统计
 */
export function getOnlineStatistics() {
  return request({
    url: '/admin/online/statistics',
    method: 'get'
  })
}
