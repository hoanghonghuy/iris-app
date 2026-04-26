import { httpClient } from './httpClient'

export const parentService = {
  async getAnalytics() {
    return await httpClient.get('/parent/analytics')
  },

  async getMyChildren() {
    return await httpClient.get('/parent/children')
  },

  async getMyFeed(params) {
    return await httpClient.get('/parent/feed', params)
  },

  async getChildPosts(studentId, params) {
    return await httpClient.get(`/parent/children/${studentId}/posts`, params)
  },

  async togglePostLike(postId) {
    return await httpClient.post(`/parent/posts/${postId}/like`)
  },

  async getPostComments(postId, params) {
    return await httpClient.get(`/parent/posts/${postId}/comments`, params)
  },

  async createPostComment(postId, data) {
    return await httpClient.post(`/parent/posts/${postId}/comments`, data)
  },

  async sharePost(postId) {
    return await httpClient.post(`/parent/posts/${postId}/share`)
  },

  async getAvailableSlots(params) {
    return await httpClient.get('/parent/appointments/slots', params)
  },

  async createAppointment(data) {
    return await httpClient.post('/parent/appointments', data)
  },

  async getAppointments(params) {
    return await httpClient.get('/parent/appointments', params)
  },

  async cancelAppointment(appointmentId, cancelReason) {
    return await httpClient.patch(`/parent/appointments/${appointmentId}/cancel`, {
      cancel_reason: cancelReason,
    })
  },
}
