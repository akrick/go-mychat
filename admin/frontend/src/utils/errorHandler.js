import { ElMessage } from 'element-plus'

/**
 * 统一错误处理
 */
export function handleError(error, defaultMessage = '操作失败') {
  console.error('Error:', error)

  if (error.response) {
    // 服务器返回错误响应
    const { status, data } = error.response

    if (status === 401) {
      ElMessage.error('登录已过期，请重新登录')
      setTimeout(() => {
        localStorage.removeItem('token')
        window.location.href = '/login'
      }, 1000)
      return
    }

    if (status === 403) {
      ElMessage.error('无权限执行此操作')
      return
    }

    if (status === 404) {
      ElMessage.error(data?.msg || '请求的资源不存在')
      return
    }

    if (status === 500) {
      ElMessage.error(data?.msg || '服务器错误，请稍后重试')
      return
    }

    // 其他错误
    ElMessage.error(data?.msg || defaultMessage)
  } else if (error.request) {
    // 请求已发送但没有收到响应
    ElMessage.error('网络错误，请检查网络连接')
  } else {
    // 其他错误
    ElMessage.error(error.message || defaultMessage)
  }
}

/**
 * 显示成功提示
 */
export function showSuccess(message = '操作成功') {
  ElMessage.success(message)
}

/**
 * 显示警告提示
 */
export function showWarning(message = '操作警告') {
  ElMessage.warning(message)
}

/**
 * 显示信息提示
 */
export function showInfo(message) {
  ElMessage.info(message)
}
