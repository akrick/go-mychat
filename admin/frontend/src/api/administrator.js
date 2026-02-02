import request from '@/utils/request'

// 管理员登录
export function adminLogin(data) {
  return request({
    url: '/api/admin2/login',
    method: 'post',
    data
  })
}

// 获取当前管理员信息
export function getAdminInfo() {
  return request({
    url: '/api/admin2/info',
    method: 'get'
  })
}

// 获取管理员列表
export function getAdministratorList(params) {
  return request({
    url: '/api/admin2/administrators',
    method: 'get',
    params
  })
}

// 创建管理员
export function createAdministrator(data) {
  return request({
    url: '/api/admin2/administrators',
    method: 'post',
    data
  })
}

// 更新管理员
export function updateAdministrator(id, data) {
  return request({
    url: `/api/admin2/administrators/${id}`,
    method: 'put',
    data
  })
}

// 删除管理员
export function deleteAdministrator(id) {
  return request({
    url: `/api/admin2/administrators/${id}`,
    method: 'delete'
  })
}

// 重置管理员密码
export function resetAdministratorPassword(id, data) {
  return request({
    url: `/api/admin2/administrators/${id}/password`,
    method: 'post',
    data
  })
}

// 切换管理员状态
export function toggleAdministratorStatus(id, data) {
  return request({
    url: `/api/admin2/administrators/${id}/status`,
    method: 'put',
    data
  })
}

// 更新个人资料
export function updateProfile(data) {
  return request({
    url: '/api/admin2/profile',
    method: 'put',
    data
  })
}

// 修改自己的密码
export function changePassword(data) {
  return request({
    url: '/api/admin2/password',
    method: 'post',
    data
  })
}

// 管理员退出登录
export function adminLogout() {
  return request({
    url: '/api/admin2/logout',
    method: 'post'
  })
}
