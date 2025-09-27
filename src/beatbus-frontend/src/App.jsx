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
import './styles/globals.css'

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
        <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50">
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
              
              {/* Protected Routes */}
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