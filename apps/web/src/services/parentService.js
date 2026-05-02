import { httpClient } from './httpClient'
import { parentPostService } from './postService'

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

  // Post interactions delegated to postService
  ...parentPostService,

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
