import { ref } from 'vue'

/**
 * 加载状态管理 Composable
 * @returns {Object}
 */
export function useLoading() {
  const loading = ref(false)
  const error = ref(null)

  /**
   * 执行异步操作并自动管理加载状态
   * @param {Function} asyncFn - 异步函数
   * @param {Object} options - 选项
   * @param {boolean} options.showError - 是否显示错误消息
   * @param {Function} options.onError - 错误回调
   * @returns {Promise}
   */
  async function executeAsync(asyncFn, options = {}) {
    const { showError = true, onError } = options

    loading.value = true
    error.value = null

    try {
      const result = await asyncFn()
      return result
    } catch (err) {
      error.value = err
      if (onError) {
        onError(err)
      }
      if (showError) {
        // 错误由调用者处理
      }
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 开始加载
   */
  function startLoading() {
    loading.value = true
    error.value = null
  }

  /**
   * 停止加载
   * @param {Error} err - 错误对象
   */
  function stopLoading(err = null) {
    loading.value = false
    error.value = err
  }

  /**
   * 重置状态
   */
  function reset() {
    loading.value = false
    error.value = null
  }

  return {
    loading,
    error,
    executeAsync,
    startLoading,
    stopLoading,
    reset
  }
}
