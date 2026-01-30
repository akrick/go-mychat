import { ElMessage } from 'element-plus'
import { REGEX, FILE_UPLOAD } from '@/constants'

/**
 * 验证手机号
 * @param {string} phone - 手机号
 * @returns {boolean}
 */
export const validatePhone = (phone) => {
  return REGEX.PHONE.test(phone)
}

/**
 * 验证邮箱
 * @param {string} email - 邮箱地址
 * @returns {boolean}
 */
export const validateEmail = (email) => {
  return REGEX.EMAIL.test(email)
}

/**
 * 验证用户名
 * @param {string} username - 用户名
 * @returns {boolean}
 */
export const validateUsername = (username) => {
  return REGEX.USERNAME.test(username)
}

/**
 * 验证密码强度
 * @param {string} password - 密码
 * @returns {boolean}
 */
export const validatePassword = (password) => {
  return REGEX.PASSWORD.test(password)
}

/**
 * 验证图片文件
 * @param {File} file - 文件对象
 * @returns {{ valid: boolean, message?: string }}
 */
export const validateImageFile = (file) => {
  if (!FILE_UPLOAD.ALLOWED_IMAGE_TYPES.includes(file.type)) {
    return {
      valid: false,
      message: '只能上传图片格式（jpeg、png、gif、webp）'
    }
  }
  if (file.size > FILE_UPLOAD.MAX_SIZE) {
    return {
      valid: false,
      message: `文件大小不能超过 ${FILE_UPLOAD.MAX_SIZE / 1024 / 1024}MB`
    }
  }
  return { valid: true }
}

/**
 * 验证表单字段
 * @param {Object} rules - 验证规则
 * @param {Object} data - 表单数据
 * @returns {Promise<{ valid: boolean, errors: Object }>}
 */
export const validateForm = async (rules, data) => {
  const errors = {}
  let valid = true

  for (const field in rules) {
    const fieldRules = rules[field]
    const value = data[field]

    for (const rule of fieldRules) {
      if (rule.required && !value) {
        errors[field] = rule.message || `${field} 是必填项`
        valid = false
        break
      }

      if (rule.pattern && value && !rule.pattern.test(value)) {
        errors[field] = rule.message || `${field} 格式不正确`
        valid = false
        break
      }

      if (rule.validator && value) {
        try {
          await rule.validator(null, value)
        } catch (error) {
          errors[field] = error.message || `${field} 验证失败`
          valid = false
          break
        }
      }
    }
  }

  return { valid, errors }
}

/**
 * 显示验证错误
 * @param {Object} errors - 错误对象
 */
export const showValidationErrors = (errors) => {
  const messages = Object.values(errors)
  if (messages.length > 0) {
    ElMessage.error(messages[0])
  }
}

/**
 * 验证确认密码
 * @param {string} password - 密码
 * @param {string} confirmPassword - 确认密码
 * @returns {boolean}
 */
export const validateConfirmPassword = (password, confirmPassword) => {
  return password === confirmPassword
}

/**
 * 验证身份证号
 * @param {string} idCard - 身份证号
 * @returns {boolean}
 */
export const validateIdCard = (idCard) => {
  return REGEX.ID_CARD.test(idCard)
}

/**
 * 验证URL
 * @param {string} url - URL地址
 * @returns {boolean}
 */
export const validateURL = (url) => {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * 获取文件扩展名
 * @param {File} file - 文件对象
 * @returns {string}
 */
export const getFileExtension = (file) => {
  return file.name.split('.').pop().toLowerCase()
}

/**
 * 格式化文件大小
 * @param {number} bytes - 字节数
 * @returns {string}
 */
export const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}
