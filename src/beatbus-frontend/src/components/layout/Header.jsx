import { Link, useNavigate } from 'react-router-dom'
import { Music, LogOut, User } from 'lucide-react'
import useAuthStore from '../../store/authStore'

const Header = () => {
  const { isAuthenticated, user, logout } = useAuthStore()
  const navigate = useNavigate()
  
  const handleLogout = () => {
    logout()
    navigate('/')
  }
  
  return (
    <header className="bg-white shadow-sm border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="flex items-center space-x-2">
            <Music className="w-8 h-8 text-primary-600" />
            <span className="text-2xl font-display font-bold text-gray-900">
              BeatBus
            </span>
          </Link>
          
          <nav className="hidden md:flex items-center space-x-8">
            <Link to="/" className="text-gray-600 hover:text-gray-900">
              Home
            </Link>
            <Link to="/join" className="text-gray-600 hover:text-gray-900">
              Join Room
            </Link>
            {isAuthenticated && (
              <Link to="/create" className="text-gray-600 hover:text-gray-900">
                Create Room
              </Link>
            )}
          </nav>
          
          <div className="flex items-center space-x-4">
            {isAuthenticated ? (
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2">
                  <User className="w-5 h-5 text-gray-500" />
                  <span className="text-gray-700">{user?.username}</span>
                </div>
                <button
                  onClick={handleLogout}
                  className="btn-ghost"
                >
                  <LogOut className="w-5 h-5" />
                </button>
              </div>
            ) : (
              <div className="flex items-center space-x-4">
                <Link to="/login" className="btn-ghost">
                  Login
                </Link>
                <Link to="/create" className="btn-primary">
                  Create Room
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}

export default Header