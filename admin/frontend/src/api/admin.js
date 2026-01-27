import request from '@/utils/request'

// 用户管理
export function getUserList(params) {
  return request({
    url: '/api/admin/users',
    method: 'get',
    params
  })
}

export function createUser(data) {
  return request({
    url: '/api/admin/users',
    method: 'post',
    data
  })
}

export function updateUser(id, data) {
  return request({
    url: `/api/admin/users/${id}`,
    method: 'put',
    data
  })
}

export function deleteUser(id) {
  return request({
    url: `/api/admin/users/${id}`,
    method: 'delete'
  })
}

export function resetUserPassword(id, data) {
  return request({
    url: `/api/admin/users/${id}/password`,
    method: 'post',
    data
  })
}

// 角色管理
export function getRoleList(params) {
  return request({
    url: '/api/admin/roles',
    method: 'get',
    params
  })
}

export function createRole(data) {
  return request({
    url: '/api/admin/roles',
    method: 'post',
    data
  })
}

export function updateRole(id, data) {
  return request({
    url: `/api/admin/roles/${id}`,
    method: 'put',
    data
  })
}

export function deleteRole(id) {
  return request({
    url: `/api/admin/roles/${id}`,
    method: 'delete'
  })
}

export function getRolePermissions(id) {
  return request({
    url: `/api/admin/roles/${id}/permissions`,
    method: 'get'
  })
}

export function assignRolePermissions(id, data) {
  return request({
    url: `/api/admin/roles/${id}/permissions`,
    method: 'put',
    data
  })
}

// 权限管理
export function getPermissionList(params) {
  return request({
    url: '/api/admin/permissions',
    method: 'get',
    params
  })
}

export function getPermissionTree() {
  return request({
    url: '/api/admin/permissions/tree',
    method: 'get'
  })
}

export function createPermission(data) {
  return request({
    url: '/api/admin/permissions',
    method: 'post',
    data
  })
}

export function updatePermission(id, data) {
  return request({
    url: `/api/admin/permissions/${id}`,
    method: 'put',
    data
  })
}

export function deletePermission(id) {
  return request({
    url: `/api/admin/permissions/${id}`,
    method: 'delete'
  })
}

// 菜单管理
export function getMenus() {
  return request({
    url: '/api/admin/menus',
    method: 'get'
  })
}

export function createMenu(data) {
  return request({
    url: '/api/admin/menus',
    method: 'post',
    data
  })
}

export function updateMenu(id, data) {
  return request({
    url: `/api/admin/menus/${id}`,
    method: 'put',
    data
  })
}

export function deleteMenu(id) {
  return request({
    url: `/api/admin/menus/${id}`,
    method: 'delete'
  })
}
