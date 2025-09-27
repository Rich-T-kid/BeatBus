import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { motion } from 'framer-motion'
import { QrCode, Hash, Camera, User } from 'lucide-react'
import { Html5QrcodeScanner } from 'html5-qrcode'
import { roomAPI } from '../services/api'
import Input from '../components/ui/Input'
import Button from '../components/ui/Button'
import Card from '../components/ui/Card'
import Modal from '../components/ui/Modal'
import toast from 'react-hot-toast'

const JoinRoom = () => {
  const navigate = useNavigate()
  const { register, handleSubmit, formState: { errors } } = useForm()
  const [isJoining, setIsJoining] = useState(false)
  const [showQRScanner, setShowQRScanner] = useState(false)
  const [joinMethod, setJoinMethod] = useState('manual') // 'manual' or 'qr'

  const onSubmit = async (data) => {
    if (!data.username?.trim()) {
      toast.error('Please enter a display name')
      return
    }

    setIsJoining(true)
    try {
      const response = await roomAPI.join(data.roomId, data.roomPassword, data.username)
      
      // Store guest info
      localStorage.setItem('guestUsername', data.username)
      localStorage.setItem('currentRoomId', data.roomId)
      
      toast.success('Joined room successfully!')
      navigate(`/room/${data.roomId}`)
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to join room')
    } finally {
      setIsJoining(false)
    }
  }

  const handleQRScan = (decodedText) => {
    try {
      // Parse QR code URL or direct room info
      const url = new URL(decodedText)
      const roomId = url.pathname.split('/').pop()
      const roomPassword = url.searchParams.get('password')
      
      if (roomId && roomPassword) {
        // Auto-fill form with QR data
        document.getElementById('roomId').value = roomId
        document.getElementById('roomPassword').value = roomPassword
        setShowQRScanner(false)
        toast.success('Room details scanned!')
      }
    } catch (error) {
      toast.error('Invalid QR code format')
    }
  }

  const QRScanner = () => {
    useState(() => {
      if (showQRScanner) {
        const scanner = new Html5QrcodeScanner(
          "qr-scanner",
          { fps: 10, qrbox: { width: 250, height: 250 } }
        )
        
        scanner.render(
          (decodedText) => {
            handleQRScan(decodedText)
            scanner.clear()
          },
          (error) => {
            console.warn('QR scan error:', error)
          }
        )
        
        return () => scanner.clear()
      }
    }, [showQRScanner])

    return <div id="qr-scanner" className="w-full" />
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 py-12 px-4">
      <div className="max-w-md mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          <Card className="p-8">
            {/* Header */}
            <div className="text-center mb-8">
              <div className="w-16 h-16 bg-secondary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <QrCode className="w-8 h-8 text-secondary-600" />
              </div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">
                Join Room
              </h1>
              <p className="text-gray-600">
                Enter room details or scan a QR code to join
              </p>
            </div>

            {/* Join Method Toggle */}
            <div className="flex rounded-lg bg-gray-100 p-1 mb-6">
              <button
                type="button"
                onClick={() => setJoinMethod('manual')}
                className={`flex-1 py-2 px-4 rounded-md text-sm font-medium transition-all ${
                  joinMethod === 'manual'
                    ? 'bg-white text-gray-900 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                <Hash className="w-4 h-4 inline mr-2" />
                Manual Entry
              </button>
              <button
                type="button"
                onClick={() => setJoinMethod('qr')}
                className={`flex-1 py-2 px-4 rounded-md text-sm font-medium transition-all ${
                  joinMethod === 'qr'
                    ? 'bg-white text-gray-900 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                <Camera className="w-4 h-4 inline mr-2" />
                QR Code
              </button>
            </div>

            {joinMethod === 'qr' ? (
              <div className="text-center">
                <Button
                  onClick={() => setShowQRScanner(true)}
                  variant="outline"
                  size="lg"
                  className="w-full mb-4"
                >
                  <Camera className="w-5 h-5 mr-2" />
                  Scan QR Code
                </Button>
                <p className="text-sm text-gray-500">
                  Scan the QR code shared by the room host
                </p>
              </div>
            ) : (
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                <Input
                  id="username"
                  label="Display Name"
                  {...register('username', { 
                    required: 'Display name is required',
                    minLength: { value: 2, message: 'Name must be at least 2 characters' },
                    maxLength: { value: 20, message: 'Name must be less than 20 characters' }
                  })}
                  error={errors.username?.message}
                  placeholder="How should others see you?"
                />

                <Input
                  id="roomId"
                  label="Room ID"
                  {...register('roomId', { 
                    required: 'Room ID is required' 
                  })}
                  error={errors.roomId?.message}
                  placeholder="Enter room ID"
                />

                <Input
                  id="roomPassword"
                  label="Room Password"
                  type="password"
                  {...register('roomPassword', { 
                    required: 'Room password is required' 
                  })}
                  error={errors.roomPassword?.message}
                  placeholder="Enter room password"
                />

                <Button
                  type="submit"
                  loading={isJoining}
                  size="lg"
                  className="w-full"
                >
                  <User className="w-5 h-5 mr-2" />
                  Join Room
                </Button>
              </form>
            )}

            {/* Help Text */}
            <div className="mt-6 pt-6 border-t border-gray-200 text-center">
              <p className="text-sm text-gray-500 mb-2">
                Don't have room details?
              </p>
              <p className="text-xs text-gray-400">
                Ask the room host to share the Room ID and password, or scan their QR code
              </p>
            </div>
          </Card>
        </motion.div>

        {/* QR Scanner Modal */}
        <Modal
          isOpen={showQRScanner}
          onClose={() => setShowQRScanner(false)}
          title="Scan QR Code"
          size="md"
        >
          <div className="text-center">
            <QRScanner />
            <p className="text-sm text-gray-500 mt-4">
              Position the QR code within the camera view
            </p>
          </div>
        </Modal>
      </div>
    </div>
  )
}

export default JoinRoom