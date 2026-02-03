import request from '@/utils/request'

/**
 * 创建支付订单
 * @param {Object} data - 支付信息
 * @param {number} data.order_id - 订单ID
 * @param {string} data.payment_method - 支付方式 wechat/alipay
 * @param {string} data.trade_type - 交易类型
 * @param {string} data.client_ip - 客户端IP
 * @param {string} data.return_url - 支付成功后跳转地址
 */
export function createPayment(data) {
  return request({
    url: '/payment/create',
    method: 'post',
    data
  })
}

/**
 * 查询支付状态
 * @param {number} id - 支付记录ID
 */
export function getPaymentStatus(id) {
  return request({
    url: `/payment/${id}`,
    method: 'get'
  })
}

/**
 * 获取用户支付记录列表
 * @param {Object} params - 查询参数
 */
export function getPaymentList(params) {
  return request({
    url: '/payment/list',
    method: 'get',
    params
  })
}

/**
 * 申请退款
 * @param {Object} data - 退款信息
 */
export function refundPayment(data) {
  return request({
    url: '/payment/refund',
    method: 'post',
    data
  })
}
