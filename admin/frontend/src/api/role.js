import request from '@/utils/request'

// 获取角色列表
export const getRoleList = (params) => {
  return request({
    url: '/api/admin/roles',
    method: 'get',
    params
  })
}

// 创建角色
export const createRole = (data) => {
  return request({
    url: '/api/admin/roles',
    method: 'post',
    data
  })
}

// 更新角色
export const updateRole = (id, data) => {
  return request({
    url: `/api/admin/roles/${id}`,
    method: 'put',
    data
  })
}

// 删除角色
export const deleteRole = (id) => {
  return request({
    url: `/api/admin/roles/${id}`,
    method: 'delete'
  })
}

// 获取角色权限
export const getRolePermissions = (id) => {
  return request({
    url: `/api/admin/roles/${id}/permissions`,
    method: 'get'
  })
}

// 分配角色权限
export const assignPermissions = (id, data) => {
  return request({
    url: `/api/admin/roles/${id}/permissions`,
    method: 'post',
    data
  })
}
