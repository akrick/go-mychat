import { ElMessage } from 'element-plus'
import { ERROR_MESSAGES, SUCCESS_MESSAGES } from '@/constants'

/**
 * 统一错误处理
 * @param {Error} error 错误对象
 * @param {string} defaultMessage 默认错误消息
 */
export function handleError(error, defaultMessage = ERROR_MESSAGES.UNKNOWN_ERROR) {
  console.error(error)

  let message = defaultMessage

  if (error.response) {
    // 服务器返回的错误
    const res = error.response.data
    message = res.msg || res.message || message
  } else if (error.message) {
    // 客户端错误
    message = error.message
  } else if (typeof error === 'string') {
    // 字符串错误
    message = error
  }

  ElMessage.error(message)
  return message
}

/**
 * 处理表单验证错误
 * @param {Error} error 验证错误
 */
export function handleValidationError(error) {
  const message = error.message || '表单验证失败'
  ElMessage.warning(message)
  return message
}

/**
 * 处理网络错误
 * @param {Error} error 网络错误
 */
export function handleNetworkError(error) {
  console.error('网络错误:', error)

  let message = ERROR_MESSAGES.NETWORK_ERROR

  if (error.code === 'ECONNABORTED') {
    message = '请求超时，请稍后重试'
  } else if (error.message && error.message.includes('Network Error')) {
    message = '网络连接失败，请检查网络'
  }

  ElMessage.error(message)
  return message
}

/**
 * 显示成功消息
 * @param {string} message 成功消息
 */
export function showSuccess(message = SUCCESS_MESSAGES.LOGIN) {
  ElMessage.success(message)
}

/**
 * 显示警告消息
 * @param {string} message 警告消息
 */
export function showWarning(message) {
  ElMessage.warning(message)
}

/**
 * 显示信息消息
 * @param {string} message 信息消息
 */
export function showInfo(message) {
  ElMessage.info(message)
}

/**
 * 显示错误消息
 * @param {string} message 错误消息
 */
export function showError(message) {
  ElMessage.error(message)
}
