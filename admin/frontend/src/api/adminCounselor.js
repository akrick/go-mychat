import request from '@/utils/request'

// 获取咨询师列表
export function getCounselorList(params) {
  return request({
    url: '/api/admin/counselors',
    method: 'get',
    params
  })
}

// 创建咨询师
export function createCounselor(data) {
  return request({
    url: '/api/admin/counselors',
    method: 'post',
    data
  })
}

// 更新咨询师
export function updateCounselor(id, data) {
  return request({
    url: `/api/admin/counselors/${id}`,
    method: 'put',
    data
  })
}

// 删除咨询师
export function deleteCounselor(id) {
  return request({
    url: `/api/admin/counselors/${id}`,
    method: 'delete'
  })
}

// 获取咨询师详情
export function getCounselorDetail(id) {
  return request({
    url: `/api/admin/counselors/${id}`,
    method: 'get'
  })
}
