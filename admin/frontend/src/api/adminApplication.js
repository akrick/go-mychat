import request from '@/utils/request'

// 获取入驻申请列表
export function getApplicationList(params) {
  return request({
    url: '/api/admin/counselor/applications',
    method: 'get',
    params
  })
}

// 审核入驻申请
export function reviewApplication(id, data) {
  return request({
    url: `/api/admin/counselor/applications/${id}/review`,
    method: 'put',
    data
  })
}

// 获取申请详情
export function getApplicationDetail(id) {
  return request({
    url: `/api/admin/counselor/applications/${id}`,
    method: 'get'
  })
}
