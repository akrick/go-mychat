import { reactive, computed, ref } from 'vue'
import { PAGINATION } from '@/constants'

/**
 * 分页管理 Composable
 * @param {Object} options - 配置选项
 * @param {number} options.pageSize - 默认每页数量
 * @param {Function} options.loadData - 加载数据的函数
 * @returns {Object}
 */
export function usePagination(options = {}) {
  const {
    pageSize = PAGINATION.DEFAULT_PAGE_SIZE,
    loadData
  } = options

  const pagination = reactive({
    page: PAGINATION.DEFAULT_PAGE,
    page_size: pageSize
  })

  const total = ref(0)
  const loading = ref(false)

  /**
   * 计算总页数
   */
  const totalPages = computed(() => {
    return Math.ceil(total.value / pagination.page_size)
  })

  /**
   * 当前页是否有数据
   */
  const hasData = computed(() => {
    return total.value > 0
  })

  /**
   * 是否是第一页
   */
  const isFirstPage = computed(() => {
    return pagination.page === 1
  })

  /**
   * 是否是最后一页
   */
  const isLastPage = computed(() => {
    return pagination.page >= totalPages.value
  })

  /**
   * 改变页码
   * @param {number} page - 页码
   */
  function changePage(page) {
    if (page < 1 || page > totalPages.value) return
    pagination.page = page
    if (loadData) {
      loadData()
    }
  }

  /**
   * 改变每页数量
   * @param {number} size - 每页数量
   */
  function changePageSize(size) {
    pagination.page_size = size
    pagination.page = 1
    if (loadData) {
      loadData()
    }
  }

  /**
   * 上一页
   */
  function prevPage() {
    if (!isFirstPage.value) {
      changePage(pagination.page - 1)
    }
  }

  /**
   * 下一页
   */
  function nextPage() {
    if (!isLastPage.value) {
      changePage(pagination.page + 1)
    }
  }

  /**
   * 重置到第一页
   */
  function resetPage() {
    pagination.page = 1
    if (loadData) {
      loadData()
    }
  }

  /**
   * 设置总数
   * @param {number} count - 总数
   */
  function setTotal(count) {
    total.value = count
  }

  /**
   * 获取分页参数
   * @returns {Object}
   */
  function getParams() {
    return {
      page: pagination.page,
      page_size: pagination.page_size
    }
  }

  return {
    pagination,
    total,
    loading,
    totalPages,
    hasData,
    isFirstPage,
    isLastPage,
    changePage,
    changePageSize,
    prevPage,
    nextPage,
    resetPage,
    setTotal,
    getParams
  }
}
