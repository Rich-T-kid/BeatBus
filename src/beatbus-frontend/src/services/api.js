import axios from 'axios'
import toast from 'react-hot-toast'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor for auth tokens
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    // Add room ID to headers if available
    const roomId = localStorage.getItem('currentRoomId')
    if (roomId) {
      config.headers['X-Room-ID'] = roomId
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    const { response } = error
    
    if (response?.status === 401) {
      // Handle unauthorized - clear tokens and redirect
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      window.location.href = '/login'
      toast.error('Session expired. Please login again.')
    } else if (response?.status >= 500) {
      toast.error('Server error. Please try again later.')
    } else if (response?.status === 429) {
      toast.error('Too many requests. Please slow down.')
    } else if (response?.data?.message) {
      toast.error(response.data.message)
    } else {
      toast.error('Something went wrong. Please try again.')
    }
    
    return Promise.reject(error)
  }
)

// API methods based on your OpenAPI spec
export const authAPI = {
  login: (credentials) => api.post('/login', credentials),
  signup: (userData) => api.post('/signUp', userData),
  refreshToken: () => api.post('/auth/refresh'),
}

export const roomAPI = {
  create: (roomData) => api.post('/rooms', roomData),
  join: (roomId, roomPassword, username) => 
    api.get(`/rooms/${roomId}`, { 
      params: { roomPassword, username } 
    }),
  update: (roomData) => api.put('/rooms', roomData),
  delete: (roomData) => api.delete('/rooms', { data: roomData }),
  getState: (roomId, roomPassword) => 
    api.get(`/rooms/${roomId}/state`, { 
      params: { roomPassword } 
    }),
}

export const queueAPI = {
  getQueue: (roomId) => api.get(`/queues/${roomId}/playlist`),
  addSong: (roomId, songData) => api.post(`/queues/${roomId}/playlist`, songData),
  reorderQueue: (roomId, newOrder) => api.put(`/queues/${roomId}/playlist`, { newOrder }),
  getNextSong: (roomId) => api.post(`/queues/${roomId}/nextSong`),
}

export const metricsAPI = {
  getRoomMetrics: (roomId) => api.get(`/metrics/${roomId}`),
  voteSong: (roomId, voteData) => api.post(`/metrics/${roomId}`, voteData),
  getHistory: (roomId) => api.get(`/metrics/${roomId}/history`),
  sendPlaylist: (roomId, playlistData) => 
    api.post(`/metrics/${roomId}/playlist/send`, playlistData),
}

export const healthAPI = {
  check: () => api.get('/health'),
}

export default api