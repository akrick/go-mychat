import request from '@/utils/request'

// 获取个人信息
export const getProfile = () => {
  return request({
    url: '/api/admin/profile',
    method: 'get'
  })
}

// 更新个人信息
export const updateProfile = (data) => {
  return request({
    url: '/api/admin/profile',
    method: 'put',
    data
  })
}

// 修改密码
export const changePassword = (data) => {
  return request({
    url: '/api/admin/user/password',
    method: 'post',
    data
  })
}

// 上传头像
export const uploadAvatar = (file) => {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/api/admin/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
