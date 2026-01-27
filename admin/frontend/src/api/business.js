import request from '@/utils/request'

// 订单管理
export function getOrderList(params) {
  return request({
    url: '/api/admin/orders',
    method: 'get',
    params
  })
}

export function updateOrderStatus(id, data) {
  return request({
    url: `/api/admin/orders/${id}/status`,
    method: 'put',
    data
  })
}

// 咨询师管理
export function getCounselorList(params) {
  return request({
    url: '/api/admin/counselors',
    method: 'get',
    params
  })
}

export function createCounselor(data) {
  return request({
    url: '/api/admin/counselors',
    method: 'post',
    data
  })
}

export function updateCounselor(id, data) {
  return request({
    url: `/api/admin/counselors/${id}`,
    method: 'put',
    data
  })
}

export function deleteCounselor(id) {
  return request({
    url: `/api/admin/counselors/${id}`,
    method: 'delete'
  })
}

// 聊天管理
export function getChatSessions(params) {
  return request({
    url: '/api/admin/chat/sessions',
    method: 'get',
    params
  })
}

export function getChatMessages(sessionId, params) {
  return request({
    url: `/api/admin/chat/sessions/${sessionId}/messages`,
    method: 'get',
    params
  })
}

export function getChatStatistics(params) {
  return request({
    url: '/api/admin/chat/statistics',
    method: 'get',
    params
  })
}

export function searchChatMessages(params) {
  return request({
    url: '/api/admin/chat/messages/search',
    method: 'get',
    params
  })
}

export function deleteChatSession(id) {
  return request({
    url: `/api/admin/chat/sessions/${id}`,
    method: 'delete'
  })
}

// 低代码平台
export function getFormList(params) {
  return request({
    url: '/api/admin/lowcode/forms',
    method: 'get',
    params
  })
}

export function saveFormDesign(data) {
  const id = data.id
  if (id) {
    return request({
      url: `/api/admin/lowcode/forms/${id}`,
      method: 'put',
      data
    })
  } else {
    return request({
      url: '/api/admin/lowcode/forms',
      method: 'post',
      data
    })
  }
}

export function deleteForm(id) {
  return request({
    url: `/api/admin/lowcode/forms/${id}`,
    method: 'delete'
  })
}

export function getFormDesign(id) {
  return request({
    url: `/api/admin/lowcode/forms/${id}`,
    method: 'get'
  })
}

export function getFormDataList(formId, params) {
  return request({
    url: `/api/admin/lowcode/forms/${formId}/data`,
    method: 'get',
    params
  })
}

export function submitFormData(formId, data) {
  return request({
    url: `/api/admin/lowcode/forms/${formId}/submit`,
    method: 'post',
    data
  })
}

export function getPageList(params) {
  return request({
    url: '/api/admin/lowcode/pages',
    method: 'get',
    params
  })
}

export function savePageDesign(data) {
  const id = data.id
  if (id) {
    return request({
      url: `/api/admin/lowcode/pages/${id}`,
      method: 'put',
      data
    })
  } else {
    return request({
      url: '/api/admin/lowcode/pages',
      method: 'post',
      data
    })
  }
}

export function deletePage(id) {
  return request({
    url: `/api/admin/lowcode/pages/${id}`,
    method: 'delete'
  })
}

export function getPageDesign(id) {
  return request({
    url: `/api/admin/lowcode/pages/${id}`,
    method: 'get'
  })
}
