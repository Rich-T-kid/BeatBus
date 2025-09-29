import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast'
import Layout from './components/layout/Layout'
import Home from './pages/Home'
import Login from './pages/Login'
import CreateRoom from './pages/CreateRoom'
import JoinRoom from './pages/JoinRoom'
import Room from './pages/Room'
import RoomMetrics from './pages/RoomMetrics'
import ProtectedRoute from './components/auth/ProtectedRoute'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <div className="min-h-screen bg-gradient-to-br from-blue-50 to-green-50">
          {window.DEMO_MODE && (
            <div className="bg-yellow-500 text-yellow-900 px-4 py-2 text-center text-sm font-medium sticky top-0 z-50">
              🎵 Demo Mode - UI Preview with Mock Data (No Backend Required)
            </div>
          )}
          
          <Toaster 
            position="top-right"
            toastOptions={{
              duration: 4000,
              style: {
                background: '#363636',
                color: '#fff',
              },
            }}
          />
          
          <Routes>
            <Route path="/" element={<Layout />}>
              <Route index element={<Home />} />
              <Route path="/login" element={<Login />} />
              <Route path="/join" element={<JoinRoom />} />
              
              <Route path="/create" element={
                <ProtectedRoute>
                  <CreateRoom />
                </ProtectedRoute>
              } />
              
              <Route path="/room/:roomId" element={<Room />} />
              <Route path="/metrics/:roomId" element={<RoomMetrics />} />
            </Route>
          </Routes>
        </div>
      </Router>
    </QueryClientProvider>
  )
}

export default App