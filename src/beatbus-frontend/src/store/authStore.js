import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authAPI } from '../services/api'
import { mockUser } from '../utils/mockData'
import toast from 'react-hot-toast'

const useAuthStore = create(
  persist(
    (set, get) => ({
      // State - initialize with demo user if demo mode
      user: window.DEMO_MODE ? mockUser : null,
      accessToken: window.DEMO_MODE ? 'demo-token' : null,
      refreshToken: window.DEMO_MODE ? 'demo-refresh' : null,
      isAuthenticated: window.DEMO_MODE ? true : false,
      isLoading: false,
      
      // Actions
      login: async (credentials) => {
        if (window.DEMO_MODE) {
          set({
            user: mockUser,
            accessToken: 'demo-token',
            refreshToken: 'demo-refresh',
            isAuthenticated: true,
          })
          toast.success('Logged in! (Demo Mode)')
          return { success: true }
        }

        set({ isLoading: true })
        try {
          const response = await authAPI.login(credentials)
          const { user, accessToken, refreshToken } = response.data
          
          set({
            user,
            accessToken,
            refreshToken,
            isAuthenticated: true,
            isLoading: false,
          })
          
          localStorage.setItem('accessToken', accessToken)
          localStorage.setItem('refreshToken', refreshToken)
          
          toast.success('Logged in successfully!')
          return { success: true }
        } catch (error) {
          set({ isLoading: false })
          return { success: false, error: error.response?.data?.message || 'Login failed' }
        }
      },
      
      signup: async (userData) => {
        if (window.DEMO_MODE) {
          toast.success('Account created! (Demo Mode)')
          return { success: true }
        }

        set({ isLoading: true })
        try {
          await authAPI.signup(userData)
          set({ isLoading: false })
          toast.success('Account created! Please login.')
          return { success: true }
        } catch (error) {
          set({ isLoading: false })
          return { success: false, error: error.response?.data?.message || 'Signup failed' }
        }
      },
      
      logout: () => {
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
        })
        
        localStorage.removeItem('accessToken')
        localStorage.removeItem('refreshToken')
        localStorage.removeItem('currentRoomId')
        
        toast.success('Logged out!')
      },
    }),
    {
      name: 'auth-storage',
    }
  )
)

export default useAuthStore