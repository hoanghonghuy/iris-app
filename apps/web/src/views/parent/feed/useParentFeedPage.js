import { computed, onMounted, ref, watch } from 'vue'
import { parentService } from '../../../services/parentService'
import { normalizeListResponse } from '../../../helpers/collectionUtils'
import { extractErrorMessage } from '../../../helpers/errorHandler'

const DEFAULT_LIMIT = 20

export function useParentFeedPage() {
  const posts = ref([])
  const children = ref([])
  const selectedChildId = ref('')
  const feedMode = ref('feed')
  const loading = ref(true)
  const errorMessage = ref('')
  const currentPage = ref(1)
  const pagination = ref({
    total: 0,
    limit: DEFAULT_LIMIT,
    offset: 0,
    has_more: false,
  })

  const currentOffset = computed(() => (currentPage.value - 1) * pagination.value.limit)
  const totalPages = computed(() => {
    const total = Number(pagination.value.total || 0)
    const limit = Number(pagination.value.limit || DEFAULT_LIMIT)
    return Math.max(1, Math.ceil(total / limit))
  })

  function buildQueryParams() {
    return {
      limit: pagination.value.limit || DEFAULT_LIMIT,
      offset: currentOffset.value,
    }
  }

  async function fetchChildren() {
    try {
      const response = await parentService.getMyChildren()
      children.value = normalizeListResponse(response)
      if (!selectedChildId.value && children.value.length > 0) {
        selectedChildId.value = children.value[0].student_id
      }
    } catch {
      children.value = []
      selectedChildId.value = ''
    }
  }

  async function fetchByMode() {
    const params = buildQueryParams()
    if (feedMode.value === 'feed') return await parentService.getMyFeed(params)
    if (!selectedChildId.value) return { data: [], pagination: { ...pagination.value, total: 0 } }
    if (feedMode.value === 'child_class') {
      return await parentService.getChildClassPosts(selectedChildId.value, params)
    }
    if (feedMode.value === 'child_student') {
      return await parentService.getChildStudentPosts(selectedChildId.value, params)
    }
    return await parentService.getChildPosts(selectedChildId.value, params)
  }

  async function fetchFeed() {
    loading.value = true
    errorMessage.value = ''

    try {
      const response = await fetchByMode()

      posts.value = normalizeListResponse(response)

      if (response?.pagination) {
        pagination.value = response.pagination
      } else {
        pagination.value = {
          total: posts.value.length,
          limit: pagination.value.limit || DEFAULT_LIMIT,
          offset: currentOffset.value,
          has_more: false,
        }
      }
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải bảng tin'
    } finally {
      loading.value = false
    }
  }

  function patchPostById(postId, patch) {
    posts.value = posts.value.map((post) =>
      post.post_id === postId ? { ...post, ...patch } : post,
    )
  }

  function handlePageChange(page) {
    currentPage.value = page
  }

  function setFeedMode(mode) {
    feedMode.value = mode
  }

  function setSelectedChild(studentId) {
    selectedChildId.value = studentId
  }

  watch(currentPage, fetchFeed)
  watch(feedMode, () => {
    if (currentPage.value !== 1) {
      currentPage.value = 1
      return
    }
    fetchFeed()
  })
  watch(selectedChildId, () => {
    if (feedMode.value === 'feed') return
    if (currentPage.value !== 1) {
      currentPage.value = 1
      return
    }
    fetchFeed()
  })

  onMounted(async () => {
    await fetchChildren()
    await fetchFeed()
  })

  return {
    posts,
    children,
    selectedChildId,
    feedMode,
    loading,
    errorMessage,
    currentPage,
    pagination,
    totalPages,
    fetchFeed,
    setFeedMode,
    setSelectedChild,
    patchPostById,
    handlePageChange,
  }
}
