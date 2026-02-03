import request from '@/utils/request'

// 获取表单列表
export const getFormList = (params) => {
  return request({
    url: '/api/admin/lowcode/forms',
    method: 'get',
    params
  })
}

// 保存表单设计
export const saveFormDesign = (data) => {
  return request({
    url: '/api/admin/lowcode/forms',
    method: 'post',
    data
  })
}

// 获取表单设计
export const getFormDesign = (id) => {
  return request({
    url: `/api/admin/lowcode/forms/${id}`,
    method: 'get'
  })
}

// 删除表单
export const deleteForm = (id) => {
  return request({
    url: `/api/admin/lowcode/forms/${id}`,
    method: 'delete'
  })
}

// 提交表单数据
export const submitFormData = (formId, data) => {
  return request({
    url: `/api/admin/lowcode/forms/${formId}/data`,
    method: 'post',
    data
  })
}

// 获取表单数据列表
export const getFormDataList = (formId, params) => {
  return request({
    url: `/api/admin/lowcode/forms/${formId}/data`,
    method: 'get',
    params
  })
}

// 删除表单数据
export const deleteFormData = (formId, dataId) => {
  return request({
    url: `/api/admin/lowcode/forms/${formId}/data/${dataId}`,
    method: 'delete'
  })
}

// 获取页面列表
export const getPageList = (params) => {
  return request({
    url: '/api/admin/lowcode/pages',
    method: 'get',
    params
  })
}

// 保存页面设计
export const savePageDesign = (data) => {
  return request({
    url: '/api/admin/lowcode/pages',
    method: 'post',
    data
  })
}

// 获取页面设计
export const getPageDesign = (id) => {
  return request({
    url: `/api/admin/lowcode/pages/${id}`,
    method: 'get'
  })
}

// 删除页面
export const deletePage = (id) => {
  return request({
    url: `/api/admin/lowcode/pages/${id}`,
    method: 'delete'
  })
}

// 预览页面
export const previewPage = (id) => {
  return request({
    url: `/api/admin/lowcode/pages/${id}/preview`,
    method: 'get'
  })
}
