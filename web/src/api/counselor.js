import request from '@/utils/request'

export function getCounselorList(params) {
  return request({
    url: '/counselor/list',
    method: 'get',
    params
  })
}

export function getCounselorDetail(id) {
  return request({
    url: `/counselor/${id}`,
    method: 'get'
  })
}

export function getCounselorReviews(id, params = {}) {
  return request({
    url: `/counselor/${id}/reviews`,
    method: 'get',
    params
  })
}

export function createCounselorReview(data) {
  return request({
    url: '/counselor/review',
    method: 'post',
    data
  })
}
