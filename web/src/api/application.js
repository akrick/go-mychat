import request from '@/utils/request'

/**
 * 提交咨询师入驻申请
 * @param {Object} data - 申请信息
 */
export function createApplication(data) {
  return request({
    url: '/counselor/application',
    method: 'post',
    data
  })
}

/**
 * 获取我的入驻申请
 */
export function getMyApplication() {
  return request({
    url: '/counselor/my-application',
    method: 'get'
  })
}

/**
 * 获取所有入驻申请（管理员）
 * @param {Object} params - 查询参数
 */
export function getAllApplications(params) {
  return request({
    url: '/counselor/applications',
    method: 'get',
    params
  })
}

/**
 * 审核入驻申请（管理员）
 * @param {Number} id - 申请ID
 * @param {Object} data - 审核信息
 */
export function reviewApplication(id, data) {
  return request({
    url: `/counselor/application/${id}/review`,
    method: 'put',
    data
  })
}

/**
 * 上传证书图片
 * @param {FormData} formData - 文件表单数据
 */
export function uploadCertificate(formData) {
  return request({
    url: '/counselor/upload-certificate',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
