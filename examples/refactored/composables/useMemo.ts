// composables/useMemo.ts
// 重构后的 Memo 组合式函数，展示最佳实践

import { ref, computed, readonly } from 'vue'
import type { MemoVO, CreateMemoDTO, UpdateMemoDTO, ListMemoParams } from '~/types'

// API 响应类型
interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

interface ListMemoResponse {
  list: MemoVO[]
  total: number
  hasNext: boolean
}

// 状态类型
interface MemoState {
  memos: MemoVO[]
  currentMemo: MemoVO | null
  loading: boolean
  error: Error | null
  pagination: {
    page: number
    size: number
    total: number
    hasNext: boolean
  }
}

/**
 * 使用 Memo 相关的状态和操作
 * @param options - 配置选项
 * @returns Memo 相关的状态和方法
 */
export function useMemo(options: { immediate?: boolean } = {}) {
  // 状态
  const state = ref<MemoState>({
    memos: [],
    currentMemo: null,
    loading: false,
    error: null,
    pagination: {
      page: 1,
      size: 10,
      total: 0,
      hasNext: false
    }
  })

  // 计算属性
  const isLoading = computed(() => state.value.loading)
  const hasError = computed(() => state.value.error !== null)
  const errorMessage = computed(() => state.value.error?.message || '')
  const isEmpty = computed(() => state.value.memos.length === 0 && !state.value.loading)
  const currentPage = computed(() => state.value.pagination.page)
  const hasMore = computed(() => state.value.pagination.hasNext)

  /**
   * 设置加载状态
   */
  function setLoading(loading: boolean) {
    state.value.loading = loading
  }

  /**
   * 设置错误状态
   */
  function setError(error: Error | null) {
    state.value.error = error
  }

  /**
   * 清除错误
   */
  function clearError() {
    state.value.error = null
  }

  /**
   * 获取 Memo 列表
   */
  async function fetchMemos(params: ListMemoParams = {}) {
    setLoading(true)
    clearError()

    try {
      const response = await $fetch<ApiResponse<ListMemoResponse>>('/api/memo/list', {
        method: 'POST',
        body: {
          page: state.value.pagination.page,
          size: state.value.pagination.size,
          ...params
        }
      })

      if (response.code !== 0) {
        throw new Error(response.message || '获取列表失败')
      }

      const { list, total, hasNext } = response.data

      // 如果是第一页，替换数据；否则追加数据
      if (state.value.pagination.page === 1) {
        state.value.memos = list
      } else {
        state.value.memos = [...state.value.memos, ...list]
      }

      state.value.pagination.total = total
      state.value.pagination.hasNext = hasNext

      return list
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  /**
   * 加载更多（下一页）
   */
  async function loadMore() {
    if (!hasMore.value || isLoading.value) return

    state.value.pagination.page++
    await fetchMemos()
  }

  /**
   * 刷新列表（回到第一页）
   */
  async function refresh() {
    state.value.pagination.page = 1
    await fetchMemos()
  }

  /**
   * 获取单个 Memo
   */
  async function fetchMemoById(id: number) {
    setLoading(true)
    clearError()

    try {
      const response = await $fetch<ApiResponse<MemoVO>>(`/api/memo/get?id=${id}`, {
        method: 'POST'
      })

      if (response.code !== 0) {
        throw new Error(response.message || '获取详情失败')
      }

      state.value.currentMemo = response.data
      return response.data
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  /**
   * 创建 Memo
   */
  async function createMemo(data: CreateMemoDTO) {
    setLoading(true)
    clearError()

    try {
      const response = await $fetch<ApiResponse<MemoVO>>('/api/memo/save', {
        method: 'POST',
        body: data
      })

      if (response.code !== 0) {
        throw new Error(response.message || '创建失败')
      }

      // 添加到列表开头
      state.value.memos.unshift(response.data)
      state.value.pagination.total++

      return response.data
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  /**
   * 更新 Memo
   */
  async function updateMemo(id: number, data: UpdateMemoDTO) {
    setLoading(true)
    clearError()

    try {
      const response = await $fetch<ApiResponse<MemoVO>>('/api/memo/save', {
        method: 'POST',
        body: { id, ...data }
      })

      if (response.code !== 0) {
        throw new Error(response.message || '更新失败')
      }

      // 更新列表中的数据
      const index = state.value.memos.findIndex(m => m.id === id)
      if (index !== -1) {
        state.value.memos[index] = { ...state.value.memos[index], ...response.data }
      }

      // 更新当前 memo
      if (state.value.currentMemo?.id === id) {
        state.value.currentMemo = { ...state.value.currentMemo, ...response.data }
      }

      return response.data
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  /**
   * 删除 Memo
   */
  async function deleteMemo(id: number) {
    setLoading(true)
    clearError()

    try {
      const response = await $fetch<ApiResponse<void>>(`/api/memo/remove?id=${id}`, {
        method: 'POST'
      })

      if (response.code !== 0) {
        throw new Error(response.message || '删除失败')
      }

      // 从列表中移除
      state.value.memos = state.value.memos.filter(m => m.id !== id)
      state.value.pagination.total--

      // 清除当前 memo
      if (state.value.currentMemo?.id === id) {
        state.value.currentMemo = null
      }
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  /**
   * 点赞 Memo
   */
  async function likeMemo(id: number) {
    try {
      const response = await $fetch<ApiResponse<void>>(`/api/memo/like?id=${id}`, {
        method: 'POST'
      })

      if (response.code !== 0) {
        throw new Error(response.message || '点赞失败')
      }

      // 更新本地状态
      const index = state.value.memos.findIndex(m => m.id === id)
      if (index !== -1) {
        state.value.memos[index] = {
          ...state.value.memos[index],
          favCount: (state.value.memos[index].favCount || 0) + 1
        }
      }
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
      throw error
    }
  }

  /**
   * 置顶/取消置顶 Memo
   */
  async function togglePin(id: number) {
    try {
      const response = await $fetch<ApiResponse<void>>(`/api/memo/setPinned?id=${id}`, {
        method: 'POST'
      })

      if (response.code !== 0) {
        throw new Error(response.message || '操作失败')
      }

      // 刷新列表以反映置顶状态变化
      await refresh()
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
      throw error
    }
  }

  /**
   * 更新 Memo 在列表中的状态
   */
  function updateMemoInList(id: number, updates: Partial<MemoVO>) {
    const index = state.value.memos.findIndex(m => m.id === id)
    if (index !== -1) {
      state.value.memos[index] = { ...state.value.memos[index], ...updates }
    }
    if (state.value.currentMemo?.id === id) {
      state.value.currentMemo = { ...state.value.currentMemo, ...updates }
    }
  }

  // 如果配置了立即加载，则自动获取数据
  if (options.immediate) {
    fetchMemos()
  }

  return {
    // 状态（只读）
    memos: readonly(computed(() => state.value.memos)),
    currentMemo: readonly(computed(() => state.value.currentMemo)),
    loading: readonly(isLoading),
    error: readonly(computed(() => state.value.error)),
    errorMessage: readonly(errorMessage),
    isEmpty: readonly(isEmpty),
    currentPage: readonly(currentPage),
    hasMore: readonly(hasMore),
    pagination: readonly(computed(() => state.value.pagination)),

    // 方法
    fetchMemos,
    loadMore,
    refresh,
    fetchMemoById,
    createMemo,
    updateMemo,
    deleteMemo,
    likeMemo,
    togglePin,
    updateMemoInList,
    clearError
  }
}

/**
 * 使用单个 Memo 的详细操作
 */
export function useMemoDetail(memoId: number) {
  const { currentMemo, loading, error, fetchMemoById, updateMemoInList } = useMemo()

  const isLoading = readonly(loading)
  const memo = readonly(currentMemo)

  async function refresh() {
    return fetchMemoById(memoId)
  }

  // 自动加载
  refresh()

  return {
    memo,
    isLoading,
    error: readonly(error),
    refresh,
    updateMemoInList
  }
}
