import request from '@/utils/request'

// 获取管理员统计数据
export function getAdminStatistics() {
  return request({
    url: '/api/admin/statistics',
    method: 'get'
  })
}

// 获取订单统计数据
export function getOrderStatistics(params) {
  return request({
    url: '/api/admin/orders/statistics',
    method: 'get',
    params
  })
}

// 获取财务统计
export function getFinanceStats(params) {
  return request({
    url: '/api/admin/finance/stats',
    method: 'get',
    params
  })
}

// 获取营收报表
export function getRevenueReport(params) {
  return request({
    url: '/api/admin/finance/revenue',
    method: 'get',
    params
  })
}

// 获取订单趋势
export function getOrderTrend(params) {
  return request({
    url: '/api/admin/stats/order',
    method: 'get',
    params
  })
}

// 获取咨询师排名
export function getCounselorRanking() {
  return request({
    url: '/api/admin/stats/counselor/ranking',
    method: 'get'
  })
}

// 获取会话统计
export function getSessionStats() {
  return request({
    url: '/api/admin/session/stats',
    method: 'get'
  })
}

// 获取在线用户
export function getOnlineUsers() {
  return request({
    url: '/api/admin/online/users',
    method: 'get'
  })
}
