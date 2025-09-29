import axios from 'axios'
import { mockRoom, mockRoomState, mockMetrics, mockHistory } from '../utils/mockData'
import toast from 'react-hot-toast'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Intercept requests in demo mode
api.interceptors.request.use(
  (config) => {
    if (window.DEMO_MODE) {
      // Return mock response instead of making real request
      return Promise.reject({ 
        config, 
        isDemoMode: true,
        mockResponse: getMockResponse(config)
      })
    }
    
    const token = localStorage.getItem('accessToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Handle demo mode responses
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.isDemoMode) {
      return Promise.resolve({ data: error.mockResponse })
    }
    
    const { response } = error
    if (response?.status === 401) {
      toast.error('Session expired')
    } else if (response?.status >= 500) {
      toast.error('Server error')
    }
    return Promise.reject(error)
  }
)

function getMockResponse(config) {
  const url = config.url
  
  if (url.includes('/rooms') && config.method === 'post') {
    return {
      roomProps: mockRoom,
      accessToken: { token: 'demo-token', expiresIn: 3600 },
      timeStamp: new Date().toISOString()
    }
  }
  
  if (url.includes('/rooms') && config.method === 'get') {
    return { success: true }
  }
  
  if (url.includes('/state')) {
    return mockRoomState
  }
  
  if (url.includes('/metrics') && url.includes('/history')) {
    return mockHistory
  }
  
  if (url.includes('/metrics')) {
    return mockMetrics
  }
  
  return { success: true }
}

export const authAPI = {
  login: (credentials) => api.post('/login', credentials),
  signup: (userData) => api.post('/signUp', userData),
}

export const roomAPI = {
  create: (roomData) => api.post('/rooms', roomData),
  join: (roomId, password, username) => 
    api.get(`/rooms/${roomId}`, { params: { roomPassword: password, username } }),
  getState: (roomId) => api.get(`/rooms/${roomId}/state`),
}

export const queueAPI = {
  addSong: (roomId, songData) => api.post(`/queues/${roomId}/playlist`, songData),
}

export const metricsAPI = {
  getRoomMetrics: (roomId) => api.get(`/metrics/${roomId}`),
  getHistory: (roomId) => api.get(`/metrics/${roomId}/history`),
  voteSong: (roomId, voteData) => api.post(`/metrics/${roomId}`, voteData),
}

export default api