import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 30000
})

request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    console.error('请求拦截器错误:', error)
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 200) {
      if (res.code === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        // 使用 router 而不是 window.location.href 以避免刷新页面
        if (window.location.pathname !== '/login') {
          window.location.href = '/login'
        }
      }
      return Promise.reject(new Error(res.msg || '请求失败'))
    }
    return res
  },
  error => {
    console.error('响应拦截器错误:', error)
    return Promise.reject(error)
  }
)

export default request
