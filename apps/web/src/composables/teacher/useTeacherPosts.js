import { ref, computed } from 'vue'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'

export function useTeacherPosts() {
  const classes = ref([])
  const selectedClassId = ref('')
  const students = ref([])
  const posts = ref([])
  const loading = ref(true)
  const loadingPosts = ref(false)
  const errorMessage = ref('')
  const pagination = ref({ total: 0, limit: 20, offset: 0, has_more: false })
  const currentPage = ref(1)

  const currentOffset = computed(() => (currentPage.value - 1) * pagination.value.limit)
  const totalPages = computed(() =>
    Math.max(1, Math.ceil((pagination.value.total || 0) / pagination.value.limit)),
  )

  async function fetchClasses() {
    loading.value = true
    errorMessage.value = ''
    try {
      const classResponse = await teacherService.getMyClasses()
      classes.value = classResponse?.data ?? []
      if (classes.value.length > 0) {
        selectedClassId.value = classes.value[0].class_id
      }
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải dữ liệu lớp học'
    } finally {
      loading.value = false
    }
  }

  async function fetchStudents() {
    if (!selectedClassId.value) {
      students.value = []
      return
    }

    try {
      const response = await teacherService.getStudentsInClass(selectedClassId.value)
      students.value = response?.data ?? []
    } catch {
      students.value = []
    }
  }

  async function fetchPosts() {
    if (!selectedClassId.value) {
      posts.value = []
      return
    }

    loadingPosts.value = true
    errorMessage.value = ''
    try {
      const response = await teacherService.getClassPosts(selectedClassId.value, {
        limit: pagination.value.limit,
        offset: currentOffset.value,
      })
      posts.value = response?.data ?? []
      if (response?.pagination) {
        pagination.value = response.pagination
      } else {
        pagination.value = {
          total: posts.value.length,
          limit: pagination.value.limit,
          offset: currentOffset.value,
          has_more: false,
        }
      }
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải bài đăng'
    } finally {
      loadingPosts.value = false
    }
  }

  function patchPostById(postId, patch) {
    posts.value = posts.value.map((post) =>
      post.post_id === postId ? { ...post, ...patch } : post,
    )
  }

  function setPage(page) {
    currentPage.value = page
  }

  function resetToFirstPage() {
    currentPage.value = 1
  }

  return {
    classes,
    selectedClassId,
    students,
    posts,
    loading,
    loadingPosts,
    errorMessage,
    pagination,
    currentPage,
    currentOffset,
    totalPages,
    fetchClasses,
    fetchStudents,
    fetchPosts,
    patchPostById,
    setPage,
    resetToFirstPage,
  }
}
