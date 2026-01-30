import request from '@/utils/request'

export function getProfile() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

export function updateProfile(data) {
  return request({
    url: '/user/profile',
    method: 'put',
    data
  })
}

export function changePassword(data) {
  return request({
    url: '/user/password',
    method: 'post',
    data
  })
}

export function uploadAvatar(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/upload/avatar',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function recharge(data) {
  return request({
    url: '/user/recharge',
    method: 'post',
    data
  })
}

export function getTransactions(params) {
  return request({
    url: '/user/transactions',
    method: 'get',
    params
  })
}
