import request from '@/utils/request'

// 获取权限树
export const getPermissionTree = () => {
  return request({
    url: '/api/admin/permissions/tree',
    method: 'get'
  })
}

// 获取权限列表
export const getPermissionList = (params) => {
  return request({
    url: '/api/admin/permissions',
    method: 'get',
    params
  })
}

// 创建权限
export const createPermission = (data) => {
  return request({
    url: '/api/admin/permissions',
    method: 'post',
    data
  })
}

// 更新权限
export const updatePermission = (id, data) => {
  return request({
    url: `/api/admin/permissions/${id}`,
    method: 'put',
    data
  })
}

// 删除权限
export const deletePermission = (id) => {
  return request({
    url: `/api/admin/permissions/${id}`,
    method: 'delete'
  })
}
