import request from '@/utils/request'

// 获取系统日志列表
export const getSystemLogs = (params) => {
  return request({
    url: '/api/admin/logs',
    method: 'get',
    params
  })
}

// 获取在线用户列表
export const getOnlineUsers = (params) => {
  return request({
    url: '/api/admin/online/users',
    method: 'get',
    params
  })
}

// 获取系统配置列表
export const getSystemConfigs = (params) => {
  return request({
    url: '/api/admin/configs',
    method: 'get',
    params
  })
}

// 创建系统配置
export const createSystemConfig = (data) => {
  return request({
    url: '/api/admin/configs',
    method: 'post',
    data
  })
}

// 更新系统配置
export const updateSystemConfig = (id, data) => {
  return request({
    url: `/api/admin/configs/${id}`,
    method: 'put',
    data
  })
}

// 批量保存配置
export const batchSaveConfigs = (configs) => {
  return request({
    url: '/api/admin/configs/batch',
    method: 'post',
    data: { configs }
  })
}

// 获取Dashboard统计数据
export const getDashboardStatistics = () => {
  return request({
    url: '/api/admin/dashboard/statistics',
    method: 'get'
  })
}

// 广播系统消息
export const broadcastMessage = (data) => {
  return request({
    url: '/api/admin/broadcast',
    method: 'post',
    data
  })
}

// 获取会话统计
export const getSessionStats = () => {
  return request({
    url: '/api/admin/session/stats',
    method: 'get'
  })
}

// 强制下线用户
export const kickOutUser = (id) => {
  return request({
    url: `/api/admin/online/users/${id}/kick`,
    method: 'post'
  })
}

// 发送消息给指定用户
export const sendToUser = (id, data) => {
  return request({
    url: `/api/admin/online/users/${id}/message`,
    method: 'post',
    data
  })
}
