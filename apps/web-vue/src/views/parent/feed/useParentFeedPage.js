import { computed, onMounted, ref, watch } from 'vue'
import { parentService } from '../../../services/parentService'
import { extractErrorMessage } from '../../../helpers/errorHandler'

const DEFAULT_LIMIT = 20

export function useParentFeedPage() {
  const posts = ref([])
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

  async function fetchFeed() {
    loading.value = true
    errorMessage.value = ''

    try {
      const response = await parentService.getMyFeed({
        limit: pagination.value.limit || DEFAULT_LIMIT,
        offset: currentOffset.value,
      })

      posts.value = Array.isArray(response?.data) ? response.data : []

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
    posts.value = posts.value.map((post) => (
      post.post_id === postId ? { ...post, ...patch } : post
    ))
  }

  function handlePageChange(page) {
    currentPage.value = page
  }

  watch(currentPage, fetchFeed)
  onMounted(fetchFeed)

  return {
    posts,
    loading,
    errorMessage,
    currentPage,
    pagination,
    totalPages,
    fetchFeed,
    patchPostById,
    handlePageChange,
  }
}
