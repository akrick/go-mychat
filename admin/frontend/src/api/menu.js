import request from '@/utils/request'

// 获取菜单树
export const getMenuTree = () => {
  return request({
    url: '/api/admin/menus/tree',
    method: 'get'
  })
}

// 获取菜单列表
export const getMenuList = () => {
  return request({
    url: '/api/admin/menus',
    method: 'get'
  })
}

// 创建菜单
export const createMenu = (data) => {
  return request({
    url: '/api/admin/menus',
    method: 'post',
    data
  })
}

// 更新菜单
export const updateMenu = (id, data) => {
  return request({
    url: `/api/admin/menus/${id}`,
    method: 'put',
    data
  })
}

// 删除菜单
export const deleteMenu = (id) => {
  return request({
    url: `/api/admin/menus/${id}`,
    method: 'delete'
  })
}
