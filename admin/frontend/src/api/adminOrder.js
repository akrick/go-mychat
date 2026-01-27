import request from '@/utils/request'

// 获取订单列表
export function getOrderList(params) {
  return request({
    url: '/api/admin/orders',
    method: 'get',
    params
  })
}

// 获取订单统计
export function getOrderStatistics() {
  return request({
    url: '/api/admin/orders/statistics',
    method: 'get'
  })
}

// 更新订单状态
export function updateOrderStatus(id, data) {
  return request({
    url: `/api/admin/orders/${id}/status`,
    method: 'put',
    data
  })
}
