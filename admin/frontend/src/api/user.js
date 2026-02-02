import request from '@/utils/request'

// 管理员登录(使用新的Administrator表)
export function login(data) {
  return request({
    url: '/api/admin2/login',
    method: 'post',
    data
  })
}

// 管理员登录 (别名)
export function adminLogin(data) {
  return request({
    url: '/api/admin2/login',
    method: 'post',
    data
  })
}

// 管理员退出
export function logout() {
  return request({
    url: '/api/admin2/logout',
    method: 'post'
  })
}

// 管理员退出 (别名)
export function adminLogout() {
  return request({
    url: '/api/admin2/logout',
    method: 'post'
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/api/admin2/info',
    method: 'get'
  })
}

// 获取用户信息 (别名)
export function getAdminInfo() {
  return request({
    url: '/api/admin2/info',
    method: 'get'
  })
}

// 获取权限
export function getPermissions() {
  return request({
    url: '/api/admin2/permissions',
    method: 'get'
  })
}

// 获取权限 (别名)
export function getAdminPermissions() {
  return request({
    url: '/api/admin2/permissions',
    method: 'get'
  })
}

// 修改密码
export function changePassword(data) {
  return request({
    url: '/api/admin2/password',
    method: 'post',
    data
  })
}

// 上传头像
export function uploadAvatar(formData) {
  return request({
    url: '/api/admin/upload/image',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

