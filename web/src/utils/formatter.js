import dayjs from 'dayjs'
import { TIME_FORMAT } from '@/constants'

/**
 * 格式化日期
 * @param {string|Date} date - 日期
 * @param {string} format - 格式
 * @returns {string}
 */
export function formatDate(date, format = TIME_FORMAT.DATE) {
  if (!date) return '-'
  return dayjs(date).format(format)
}

/**
 * 格式化日期时间
 * @param {string|Date} datetime - 日期时间
 * @param {string} format - 格式
 * @returns {string}
 */
export function formatDateTime(datetime, format = TIME_FORMAT.DATETIME) {
  if (!datetime) return '-'
  return dayjs(datetime).format(format)
}

/**
 * 格式化时间
 * @param {string|Date} time - 时间
 * @param {string} format - 格式
 * @returns {string}
 */
export function formatTime(time, format = TIME_FORMAT.TIME) {
  if (!time) return '-'
  return dayjs(time).format(format)
}

/**
 * 格式化金额
 * @param {number} amount - 金额
 * @param {number} decimals - 小数位数
 * @param {string} currency - 货币符号
 * @returns {string}
 */
export function formatCurrency(amount, decimals = 2, currency = '¥') {
  if (amount === null || amount === undefined) return '-'
  const fixed = Number(amount).toFixed(decimals)
  return `${currency}${fixed}`
}

/**
 * 格式化数字，添加千位分隔符
 * @param {number} num - 数字
 * @returns {string}
 */
export function formatNumber(num) {
  if (num === null || num === undefined) return '-'
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

/**
 * 格式化百分比
 * @param {number} value - 值
 * @param {number} decimals - 小数位数
 * @returns {string}
 */
export function formatPercent(value, decimals = 2) {
  if (value === null || value === undefined) return '-'
  return `${(value * 100).toFixed(decimals)}%`
}

/**
 * 格式化文件大小
 * @param {number} bytes - 字节数
 * @returns {string}
 */
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

/**
 * 格式化秒数为时间
 * @param {number} seconds - 秒数
 * @returns {string}
 */
export function formatSeconds(seconds) {
  if (!seconds) return '00:00'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)

  if (h > 0) {
    return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  }
  return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
}

/**
 * 格式化时长为可读文本
 * @param {number} minutes - 分钟数
 * @returns {string}
 */
export function formatDuration(minutes) {
  if (!minutes) return '-'
  const h = Math.floor(minutes / 60)
  const m = minutes % 60

  if (h > 0) {
    return `${h}小时${m > 0 ? m + '分钟' : ''}`
  }
  return `${m}分钟`
}

/**
 * 截断文本
 * @param {string} text - 文本
 * @param {number} length - 最大长度
 * @param {string} suffix - 后缀
 * @returns {string}
 */
export function truncateText(text, length = 100, suffix = '...') {
  if (!text) return ''
  if (text.length <= length) return text
  return text.substring(0, length) + suffix
}

/**
 * 高亮关键词
 * @param {string} text - 文本
 * @param {string} keyword - 关键词
 * @param {string} className - 类名
 * @returns {string}
 */
export function highlightKeyword(text, keyword, className = 'highlight') {
  if (!text || !keyword) return text
  const regex = new RegExp(`(${keyword})`, 'gi')
  return text.replace(regex, `<span class="${className}">$1</span>`)
}

/**
 * 格式化相对时间
 * @param {string|Date} date - 日期
 * @returns {string}
 */
export function formatRelativeTime(date) {
  if (!date) return '-'
  const now = dayjs()
  const target = dayjs(date)
  const diff = now.diff(target, 'second')

  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  if (diff < 2592000) return `${Math.floor(diff / 86400)}天前`
  if (diff < 31536000) return `${Math.floor(diff / 2592000)}个月前`
  return `${Math.floor(diff / 31536000)}年前`
}

/**
 * 隐藏手机号中间4位
 * @param {string} phone - 手机号
 * @returns {string}
 */
export function maskPhone(phone) {
  if (!phone || phone.length !== 11) return phone
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

/**
 * 隐藏身份证号中间部分
 * @param {string} idCard - 身份证号
 * @returns {string}
 */
export function maskIdCard(idCard) {
  if (!idCard || idCard.length < 8) return idCard
  const len = idCard.length
  const start = idCard.substring(0, 4)
  const end = idCard.substring(len - 4)
  const mask = '*'.repeat(len - 8)
  return start + mask + end
}

/**
 * 隐藏邮箱
 * @param {string} email - 邮箱
 * @returns {string}
 */
export function maskEmail(email) {
  if (!email || !email.includes('@')) return email
  const [username, domain] = email.split('@')
  if (username.length <= 2) return email
  const start = username.substring(0, 2)
  const end = username.substring(username.length - 1)
  const mask = '*'.repeat(username.length - 3)
  return `${start}${mask}${end}@${domain}`
}

/**
 * 生成随机字符串
 * @param {number} length - 长度
 * @returns {string}
 */
export function randomString(length = 8) {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

/**
 * 生成随机数字
 * @param {number} length - 长度
 * @returns {string}
 */
export function randomNumber(length = 6) {
  let result = ''
  for (let i = 0; i < length; i++) {
    result += Math.floor(Math.random() * 10)
  }
  return result
}

/**
 * 复制到剪贴板
 * @param {string} text - 文本
 * @returns {Promise<boolean>}
 */
export async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch (err) {
    console.error('复制失败:', err)
    return false
  }
}

/**
 * 下载文件
 * @param {string} url - 文件URL
 * @param {string} filename - 文件名
 */
export function downloadFile(url, filename) {
  const link = document.createElement('a')
  link.href = url
  link.download = filename || 'download'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

/**
 * 打开新窗口
 * @param {string} url - URL
 * @param {string} target - 目标
 */
export function openWindow(url, target = '_blank') {
  window.open(url, target)
}
