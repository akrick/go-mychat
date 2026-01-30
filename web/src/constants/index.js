// 路由名称常量
export const ROUTE_NAMES = {
  HOME: 'Home',
  LOGIN: 'Login',
  REGISTER: 'Register',
  SERVICES: 'Services',
  ABOUT: 'About',
  COUNSELORS: 'Counselors',
  COUNSELOR_DETAIL: 'CounselorDetail',
  ORDERS: 'Orders',
  CHAT: 'Chat',
  PROFILE: 'Profile',
  FAQ: 'FAQ',
  PRIVACY: 'Privacy',
  TERMS: 'Terms',
  GUIDE: 'Guide',
  NOT_FOUND: 'NotFound'
}

// 订单状态
export const ORDER_STATUS = {
  PENDING: 'pending',
  PAID: 'paid',
  IN_PROGRESS: 'in_progress',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded'
}

// 订单状态显示文本
export const ORDER_STATUS_TEXT = {
  [ORDER_STATUS.PENDING]: '待支付',
  [ORDER_STATUS.PAID]: '已支付',
  [ORDER_STATUS.IN_PROGRESS]: '进行中',
  [ORDER_STATUS.COMPLETED]: '已完成',
  [ORDER_STATUS.CANCELLED]: '已取消',
  [ORDER_STATUS.REFUNDED]: '已退款'
}

// 订单状态颜色
export const ORDER_STATUS_COLOR = {
  [ORDER_STATUS.PENDING]: 'warning',
  [ORDER_STATUS.PAID]: 'success',
  [ORDER_STATUS.IN_PROGRESS]: 'primary',
  [ORDER_STATUS.COMPLETED]: 'success',
  [ORDER_STATUS.CANCELLED]: 'info',
  [ORDER_STATUS.REFUNDED]: 'danger'
}

// 会话时长选项
export const SESSION_DURATION = {
  MINUTES_30: { value: 30, label: '30分钟', desc: '适合初步沟通' },
  MINUTES_60: { value: 60, label: '60分钟', desc: '标准咨询时长' },
  MINUTES_90: { value: 90, label: '90分钟', desc: '深度咨询' },
  MINUTES_120: { value: 120, label: '120分钟', desc: '综合治疗' }
}

// 分页默认值
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_PAGE_SIZE: 10,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100]
}

// 首页展示数量
export const HOME_COUNSELOR_COUNT = 4

// 文件上传限制
export const FILE_UPLOAD = {
  MAX_SIZE: 5 * 1024 * 1024, // 5MB
  ALLOWED_IMAGE_TYPES: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
  ALLOWED_FILE_TYPES: ['image/*', 'application/pdf', 'application/msword']
}

// 正则表达式
export const REGEX = {
  PHONE: /^1[3-9]\d{9}$/,
  EMAIL: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
  USERNAME: /^[a-zA-Z0-9_]{3,50}$/,
  PASSWORD: /^.{6,}$/,
  ID_CARD: /^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$/
}

// 消息类型
export const MESSAGE_TYPE = {
  TEXT: 'text',
  IMAGE: 'image',
  FILE: 'file',
  VOICE: 'voice',
  VIDEO: 'video'
}

// 发送者类型
export const SENDER_TYPE = {
  USER: 'user',
  COUNSELOR: 'counselor',
  SYSTEM: 'system'
}

// 咨询师评分
export const RATING = {
  MIN: 1,
  MAX: 5,
  DEFAULT: 0
}

// 时间格式
export const TIME_FORMAT = {
  DATE: 'YYYY-MM-DD',
  DATETIME: 'YYYY-MM-DD HH:mm',
  DATETIME_SECONDS: 'YYYY-MM-DD HH:mm:ss',
  TIME: 'HH:mm',
  MONTH: 'YYYY年MM月'
}

// 存储键名
export const STORAGE_KEYS = {
  TOKEN: 'token',
  USER_INFO: 'userInfo',
  THEME: 'theme',
  LANGUAGE: 'language'
}

// 错误提示
export const ERROR_MESSAGES = {
  NETWORK_ERROR: '网络错误，请检查网络连接',
  SERVER_ERROR: '服务器错误，请稍后重试',
  UNAUTHORIZED: '未授权，请重新登录',
  FORBIDDEN: '没有权限访问',
  NOT_FOUND: '资源不存在',
  VALIDATION_ERROR: '输入数据有误',
  UNKNOWN_ERROR: '未知错误'
}

// 成功提示
export const SUCCESS_MESSAGES = {
  LOGIN: '登录成功',
  REGISTER: '注册成功',
  LOGOUT: '退出成功',
  UPDATE_PROFILE: '个人信息更新成功',
  CHANGE_PASSWORD: '密码修改成功',
  UPLOAD_AVATAR: '头像上传成功',
  CREATE_ORDER: '订单创建成功',
  CANCEL_ORDER: '订单已取消',
  SEND_MESSAGE: '消息发送成功',
  SUBMIT_REVIEW: '评价提交成功'
}

// 咨询时长（秒）
export const CONSULTATION_DURATION = {
  MINUTES_30: 30 * 60,
  MINUTES_60: 60 * 60,
  MINUTES_90: 90 * 60,
  MINUTES_120: 120 * 60
}

// 咨询倒计时警告时间（秒）
export const COUNTDOWN_WARNING = {
  WARNING: 300, // 5分钟
  DANGER: 60 // 1分钟
}
