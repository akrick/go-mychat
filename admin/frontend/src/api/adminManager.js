import request from '@/utils/request'

// 获取管理员列表
export function getAdminList(params) {
  return request({
    url: '/api/admin/managers',
    method: 'get',
    params
  })
}

// 创建管理员
export function createAdmin(data) {
  return request({
    url: '/api/admin/managers',
    method: 'post',
    data
  })
}

// 更新管理员
export function updateAdmin(id, data) {
  return request({
    url: `/api/admin/managers/${id}`,
    method: 'put',
    data
  })
}

// 删除管理员
export function deleteAdmin(id) {
  return request({
    url: `/api/admin/managers/${id}`,
    method: 'delete'
  })
}

// 重置管理员密码
export function resetAdminPassword(id, data) {
  return request({
    url: `/api/admin/managers/${id}/password`,
    method: 'post',
    data
  })
}

// 禁用/启用管理员
export function toggleAdminStatus(id, data) {
  return request({
    url: `/api/admin/managers/${id}/status`,
    method: 'put',
    data
  })
}
