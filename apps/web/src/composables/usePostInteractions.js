import { ref } from 'vue'
import { teacherService } from '@/services/teacherService'
import { parentService } from '@/services/parentService'
import { extractErrorMessage } from '@/helpers/errorHandler'

const serviceMap = {
  teacher: teacherService,
  parent: parentService,
}

/**
 * Composable for post interactions (like, comment, share)
 * Handles all post interaction logic for both teacher and parent roles
 * 
 * @param {string} audience - Role audience ('teacher' or 'parent')
 * @returns {Object} Post interaction methods and state
 */
export function usePostInteractions(audience) {
  const processing = ref(false)
  const loadingComments = ref(false)
  const submittingComment = ref(false)
  const error = ref('')

  const service = serviceMap[audience]

  if (!service) {
    console.error(`Invalid audience: ${audience}. Must be 'teacher' or 'parent'`)
  }

  /**
   * Toggle like status for a post
   * @param {string} postId - Post ID
   * @returns {Promise<Object>} Updated like data
   */
  async function toggleLike(postId) {
    processing.value = true
    error.value = ''
    try {
      const response = await service.togglePostLike(postId)
      const payload = response?.data ?? response
      return payload
    } catch (err) {
      error.value = 'Không thể cập nhật lượt thích. Vui lòng thử lại.'
      throw err
    } finally {
      processing.value = false
    }
  }

  /**
   * Load comments for a post
   * @param {string} postId - Post ID
   * @param {Object} params - Query parameters (limit, offset)
   * @returns {Promise<Array>} List of comments
   */
  async function loadComments(postId, params = { limit: 50, offset: 0 }) {
    loadingComments.value = true
    error.value = ''
    try {
      const response = await service.getPostComments(postId, params)
      const comments = response?.data ?? []
      return comments
    } catch (err) {
      error.value = 'Không thể tải bình luận. Vui lòng thử lại.'
      throw err
    } finally {
      loadingComments.value = false
    }
  }

  /**
   * Create a new comment on a post
   * @param {string} postId - Post ID
   * @param {string} content - Comment content
   * @returns {Promise<Object>} Created comment data
   */
  async function createComment(postId, content) {
    if (!content || !content.trim()) {
      error.value = 'Nội dung bình luận không được để trống'
      return null
    }

    submittingComment.value = true
    error.value = ''
    try {
      const response = await service.createPostComment(postId, { content: content.trim() })
      const payload = response?.data ?? response
      return payload
    } catch (err) {
      error.value = extractErrorMessage(err) || 'Không thể gửi bình luận. Vui lòng thử lại.'
      throw err
    } finally {
      submittingComment.value = false
    }
  }

  /**
   * Share a post
   * @param {string} postId - Post ID
   * @returns {Promise<Object>} Updated share data
   */
  async function share(postId) {
    processing.value = true
    error.value = ''
    try {
      const response = await service.sharePost(postId)
      const payload = response?.data ?? response
      return payload
    } catch (err) {
      error.value = 'Không thể chia sẻ bài viết. Vui lòng thử lại.'
      throw err
    } finally {
      processing.value = false
    }
  }

  /**
   * Clear error message
   */
  function clearError() {
    error.value = ''
  }

  return {
    // State
    processing,
    loadingComments,
    submittingComment,
    error,

    // Methods
    toggleLike,
    loadComments,
    createComment,
    share,
    clearError,
  }
}
