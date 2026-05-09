import { httpClient } from './httpClient'

/**
 * Factory function to create post service for a specific role
 * Eliminates code duplication between teacherService and parentService
 * 
 * @param {string} rolePrefix - Role prefix for API endpoints ('teacher' or 'parent')
 * @returns {Object} Post service methods
 */
export function createPostService(rolePrefix) {
  const basePath = `/${rolePrefix}/posts`

  return {
    /**
     * Toggle like status for a post
     * @param {string} postId - Post ID
     * @returns {Promise<Object>} Response with liked_by_me and like_count
     */
    async togglePostLike(postId) {
      return await httpClient.post(`${basePath}/${postId}/like`)
    },

    /**
     * Get comments for a post
     * @param {string} postId - Post ID
     * @param {Object} params - Query parameters (limit, offset)
     * @returns {Promise<Object>} Response with comments array
     */
    async getPostComments(postId, params) {
      return await httpClient.get(`${basePath}/${postId}/comments`, params)
    },

    /**
     * Create a new comment on a post
     * @param {string} postId - Post ID
     * @param {Object} data - Comment data { content }
     * @returns {Promise<Object>} Response with created comment
     */
    async createPostComment(postId, data) {
      return await httpClient.post(`${basePath}/${postId}/comments`, data)
    },

    /**
     * Share a post
     * @param {string} postId - Post ID
     * @returns {Promise<Object>} Response with share_count
     */
    async sharePost(postId) {
      return await httpClient.post(`${basePath}/${postId}/share`)
    },
  }
}

/**
 * Shared post service instance for teacher role
 */
export const teacherPostService = {
  ...createPostService('teacher'),
  async updatePost(postId, data) {
    return await httpClient.put(`/teacher/posts/${postId}`, data)
  },
  async deletePost(postId) {
    return await httpClient.del(`/teacher/posts/${postId}`)
  },
}

/**
 * Shared post service instance for parent role
 */
export const parentPostService = createPostService('parent')
