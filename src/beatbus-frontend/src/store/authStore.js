import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authAPI } from '../services/api'
import toast from 'react-hot-toast'

const useAuthStore = create(
  persist(
    (set, get) => ({
      // State
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,
      
      // Actions
      login: async (credentials) => {
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
          
          // Store tokens in localStorage for axios interceptor
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
        set({ isLoading: true })
        try {
          const response = await authAPI.signup(userData)
          
          set({ isLoading: false })
          toast.success('Account created successfully! Please login.')
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
        
        // Clear localStorage
        localStorage.removeItem('accessToken')
        localStorage.removeItem('refreshToken')
        localStorage.removeItem('currentRoomId')
        
        toast.success('Logged out successfully!')
      },
      
      refreshAuthToken: async () => {
        try {
          const response = await authAPI.refreshToken()
          const { accessToken } = response.data
          
          set({ accessToken })
          localStorage.setItem('accessToken', accessToken)
          
          return { success: true }
        } catch (error) {
          // If refresh fails, logout user
          get().logout()
          return { success: false }
        }
      },
      
      // Initialize auth state from localStorage
      initializeAuth: () => {
        const accessToken = localStorage.getItem('accessToken')
        const refreshToken = localStorage.getItem('refreshToken')
        
        if (accessToken && refreshToken) {
          set({
            accessToken,
            refreshToken,
            isAuthenticated: true,
          })
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)

export default useAuthStore