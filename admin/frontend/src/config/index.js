/**
 * 全局配置文件
 */

// 应用配置
export const appConfig = {
  appName: 'MyChat 管理后台',
  appVersion: '1.1.0',
  apiBaseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  apiTimeout: 10000
}

// 分页配置
export const paginationConfig = {
  defaultPageSize: 10,
  pageSizes: [10, 20, 50, 100],
  layout: 'total, sizes, prev, pager, next, jumper'
}

// 表格配置
export const tableConfig = {
  stripe: true,
  border: true,
  size: 'default',
  showHeader: true,
  highlightCurrentRow: true
}

// 上传配置
export const uploadConfig = {
  maxSize: 10 * 1024 * 1024, // 10MB
  allowedTypes: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
  allowedExtensions: ['.jpg', '.jpeg', '.png', '.gif', '.webp']
}

// 主题配置
export const themeConfig = {
  primaryColor: '#667eea',
  successColor: '#67c23a',
  warningColor: '#e6a23c',
  dangerColor: '#f56c6c',
  infoColor: '#909399'
}

// 路由配置
export const routeConfig = {
  homePath: '/',
  loginPath: '/login',
  error404Path: '/404'
}

// 本地存储键名
export const storageKeys = {
  TOKEN: 'token',
  USER_INFO: 'user_info',
  ROLES: 'roles',
  PERMISSIONS: 'permissions',
  LANGUAGE: 'language',
  THEME: 'theme',
  SIDEBAR_COLLAPSED: 'sidebar_collapsed'
}

// 日期格式
export const dateFormat = {
  date: 'YYYY-MM-DD',
  datetime: 'YYYY-MM-DD HH:mm:ss',
  time: 'HH:mm:ss',
  month: 'YYYY-MM',
  year: 'YYYY'
}

// 正则表达式
export const regExp = {
  phone: /^1[3-9]\d{9}$/,
  email: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
  password: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,20}$/,
  url: /^https?:\/\/.+$/,
  number: /^\d+$/
}

// 错误码映射
export const errorCodeMap = {
  400: '请求参数错误',
  401: '未授权，请登录',
  403: '拒绝访问',
  404: '请求资源不存在',
  405: '请求方法不允许',
  408: '请求超时',
  500: '服务器内部错误',
  502: '网关错误',
  503: '服务不可用',
  504: '网关超时'
}

// 默认数据
export const defaultData = {
  avatar: 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxMDAiIGhlaWdodD0iMTAwIj48cmVjdCB3aWR0aD0iNjUiIGhlaWdodD0iNjUiIHJ4PSI1MDAiIHJ5PSI3MCIgZmlsbD0iI2QxMjZmZiIi8+PC9zdmc+',
  pageSize: 10,
  currentPage: 1
}

// 文件类型映射
export const fileTypeMap = {
  image: ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp'],
  document: ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'],
  video: ['mp4', 'avi', 'mov', 'wmv'],
  audio: ['mp3', 'wav', 'flac'],
  archive: ['zip', 'rar', '7z', 'tar', 'gz']
}

// 状态码映射
export const statusMap = {
  success: 200,
  created: 201,
  noContent: 204,
  badRequest: 400,
  unauthorized: 401,
  forbidden: 403,
  notFound: 404,
  serverError: 500
}

// 常用枚举
export const enums = {
  // 用户状态
  userStatus: {
    active: 1,
    inactive: 0,
    banned: -1
  },
  // 订单状态
  orderStatus: {
    pending: 0,
    paid: 1,
    completed: 2,
    cancelled: 3,
    refunded: 4
  },
  // 提现状态
  withdrawStatus: {
    pending: 0,
    approved: 1,
    rejected: 2,
    completed: 3
  },
  // 角色类型
  roleType: {
    admin: 'admin',
    counselor: 'counselor',
    user: 'user'
  }
}

// 工具函数
export const utils = {
  // 格式化文件大小
  formatFileSize: (bytes) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
  },

  // 格式化数字
  formatNumber: (num) => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + '万'
    }
    return num.toLocaleString()
  },

  // 防抖
  debounce: (fn, delay) => {
    let timer = null
    return function (...args) {
      clearTimeout(timer)
      timer = setTimeout(() => fn.apply(this, args), delay)
    }
  },

  // 节流
  throttle: (fn, delay) => {
    let last = 0
    return function (...args) {
      const now = Date.now()
      if (now - last >= delay) {
        last = now
        fn.apply(this, args)
      }
    }
  },

  // 深拷贝
  deepClone: (obj) => {
    if (obj === null || typeof obj !== 'object') return obj
    if (obj instanceof Date) return new Date(obj.getTime())
    if (obj instanceof Array) return obj.map(item => utils.deepClone(item))
    if (obj instanceof Object) {
      const clonedObj = {}
      for (const key in obj) {
        if (obj.hasOwnProperty(key)) {
          clonedObj[key] = utils.deepClone(obj[key])
        }
      }
      return clonedObj
    }
  }
}

export default {
  appConfig,
  paginationConfig,
  tableConfig,
  uploadConfig,
  themeConfig,
  routeConfig,
  storageKeys,
  dateFormat,
  regExp,
  errorCodeMap,
  defaultData,
  fileTypeMap,
  statusMap,
  enums,
  utils
}
