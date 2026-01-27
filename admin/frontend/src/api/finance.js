import request from '@/utils/request'

// 提现管理
export function getPendingWithdraws(params) {
  return request({
    url: '/api/admin/withdraws/pending',
    method: 'get',
    params
  })
}

export function getWithdrawList(params) {
  return request({
    url: '/api/admin/withdraws',
    method: 'get',
    params
  })
}

export function approveWithdraw(id, data) {
  return request({
    url: `/api/admin/withdraw/${id}/approve`,
    method: 'post',
    data
  })
}

// 财务统计
export function getFinanceStats(params) {
  return request({
    url: '/api/admin/finance/stats',
    method: 'get',
    params
  })
}

export function getRevenueReport(params) {
  return request({
    url: '/api/admin/finance/revenue',
    method: 'get',
    params
  })
}

// 财务报表
export function getFinanceReports(params) {
  return request({
    url: '/api/admin/finance/reports',
    method: 'get',
    params
  })
}
