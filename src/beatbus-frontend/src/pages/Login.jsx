import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { motion } from 'framer-motion'
import { Music, Eye, EyeOff } from 'lucide-react'
import useAuthStore from '../store/authStore'
import Input from '../components/ui/Input'
import Button from '../components/ui/Button'
import Card from '../components/ui/Card'

const Login = () => {
  const navigate = useNavigate()
  const { login, isLoading } = useAuthStore()
  const { register, handleSubmit, formState: { errors } } = useForm()
  const [showPassword, setShowPassword] = useState(false)
  const [isSignup, setIsSignup] = useState(false)

  const onSubmit = async (data) => {
    const result = await login(data)
    if (result.success) {
      navigate('/')
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 flex items-center justify-center px-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="w-full max-w-md"
      >
        <Card className="p-8">
          {/* Logo */}
          <div className="flex items-center justify-center mb-8">
            <Music className="w-8 h-8 text-primary-600 mr-2" />
            <h1 className="text-2xl font-display font-bold text-gray-900">
              BeatBus
            </h1>
          </div>

          {/* Title */}
          <div className="text-center mb-8">
            <h2 className="text-2xl font-bold text-gray-900 mb-2">
              Welcome Back
            </h2>
            <p className="text-gray-600">
              Sign in to create and manage your music rooms
            </p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <Input
              label="Username"
              type="text"
              {...register('username', { 
                required: 'Username is required',
                minLength: { value: 3, message: 'Username must be at least 3 characters' }
              })}
              error={errors.username?.message}
              placeholder="Enter your username"
            />

            <div className="relative">
              <Input
                label="Password"
                type={showPassword ? 'text' : 'password'}
                {...register('password', { 
                  required: 'Password is required',
                  minLength: { value: 6, message: 'Password must be at least 6 characters' }
                })}
                error={errors.password?.message}
                placeholder="Enter your password"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-8 text-gray-400 hover:text-gray-600"
              >
                {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>

            <Button
              type="submit"
              className="w-full"
              loading={isLoading}
              size="lg"
            >
              Sign In
            </Button>
          </form>

          {/* Signup Link */}
          <div className="mt-6 text-center">
            <p className="text-gray-600">
              Don't have an account?{' '}
              <button
                onClick={() => setIsSignup(true)}
                className="text-primary-600 hover:text-primary-700 font-medium"
              >
                Sign up
              </button>
            </p>
          </div>

          {/* Guest Join */}
          <div className="mt-6 pt-6 border-t border-gray-200">
            <div className="text-center">
              <p className="text-gray-600 mb-4">
                Just want to join a room?
              </p>
              <Link
                to="/join"
                className="btn-outline w-full"
              >
                Join as Guest
              </Link>
            </div>
          </div>
        </Card>

        {/* Signup Modal would go here if isSignup is true */}
      </motion.div>
    </div>
  )
}

export default Login