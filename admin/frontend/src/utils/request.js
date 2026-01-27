import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import router from '@/router'

const service = axios.create({
  baseURL: '',
  timeout: 10000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    console.log('响应数据:', res)
    
    // 检查响应格式
    if (res.code !== undefined) {
      if (res.code !== 200) {
        ElMessage.error(res.msg || '请求失败')

        // 401: 未授权
        if (res.code === 401) {
          ElMessageBox.confirm('登录状态已过期，请重新登录', '提示', {
            confirmButtonText: '重新登录',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            localStorage.removeItem('token')
            router.push('/login')
          })
        }
        return Promise.reject(new Error(res.msg || '请求失败'))
      }
      return res
    }
    // 如果没有code字段，直接返回（可能不是标准格式）
    return res
  },
  error => {
    console.error('请求错误:', error)
    
    // 处理HTTP状态码错误
    if (error.response) {
      const status = error.response.status
      const data = error.response.data
      
      if (status === 401) {
        ElMessage.error(data?.msg || '用户名或密码错误')
      } else if (status === 403) {
        ElMessage.error(data?.msg || '无权限访问')
      } else if (status === 500) {
        ElMessage.error(data?.msg || '服务器错误')
      } else {
        ElMessage.error(data?.msg || error.message || '网络错误')
      }
    } else {
      ElMessage.error(error.message || '网络错误')
    }
    
    return Promise.reject(error)
  }
)

export default service
