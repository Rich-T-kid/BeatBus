import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { motion } from 'framer-motion'
import { Music, Users, Clock, Lock, Globe } from 'lucide-react'
import { QRCodeSVG } from 'qrcode.react'
import useAuthStore from '../store/authStore'
import { roomAPI } from '../services/api'
import Input from '../components/ui/Input'
import Button from '../components/ui/Button'
import Card from '../components/ui/Card'
import Modal from '../components/ui/Modal'
import toast from 'react-hot-toast'

const CreateRoom = () => {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const { register, handleSubmit, watch, formState: { errors } } = useForm({
    defaultValues: {
      hostUsername: user?.username,
      roomName: '',
      lifeTime: 60,
      maxUsers: 50,
      isPublic: true
    }
  })
  
  const [isCreating, setIsCreating] = useState(false)
  const [roomCreated, setRoomCreated] = useState(null)
  const [showQR, setShowQR] = useState(false)
  
  const watchIsPublic = watch('isPublic')

  const onSubmit = async (data) => {
    setIsCreating(true)
    try {
      const response = await roomAPI.create(data)
      const { roomProps, accessToken } = response.data
      
      // Store room info
      localStorage.setItem('currentRoomId', roomProps.roomId)
      localStorage.setItem('accessToken', accessToken.token)
      
      setRoomCreated(roomProps)
      toast.success('Room created successfully!')
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to create room')
    } finally {
      setIsCreating(false)
    }
  }

  const handleJoinRoom = () => {
    if (roomCreated) {
      navigate(`/room/${roomCreated.roomId}`)
    }
  }

  const getRoomJoinUrl = () => {
    if (!roomCreated) return ''
    return `${window.location.origin}/room/${roomCreated.roomId}?password=${roomCreated.roomPassword}`
  }

  if (roomCreated) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 flex items-center justify-center px-4">
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.6 }}
          className="w-full max-w-2xl"
        >
          <Card className="p-8 text-center">
            <div className="mb-6">
              <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Music className="w-8 h-8 text-green-600" />
              </div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">
                Room Created!
              </h1>
              <p className="text-gray-600">
                Your music room "{roomCreated.roomName}" is ready
              </p>
            </div>

            <div className="bg-gray-50 rounded-xl p-6 mb-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div>
                  <span className="font-medium text-gray-900">Room ID:</span>
                  <p className="text-primary-600 font-mono">{roomCreated.roomId}</p>
                </div>
                <div>
                  <span className="font-medium text-gray-900">Password:</span>
                  <p className="text-primary-600 font-mono">{roomCreated.roomPassword}</p>
                </div>
                <div>
                  <span className="font-medium text-gray-900">Max Users:</span>
                  <p className="text-gray-600">{roomCreated.maxUsers}</p>
                </div>
                <div>
                  <span className="font-medium text-gray-900">Type:</span>
                  <p className="text-gray-600">
                    {roomCreated.isPublic ? 'Public' : 'Private'}
                  </p>
                </div>
              </div>
            </div>

            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-6">
              <Button
                onClick={handleJoinRoom}
                size="lg"
                className="flex-1 sm:flex-none"
              >
                Enter Room
              </Button>
              
              <Button
                onClick={() => setShowQR(true)}
                variant="outline"
                size="lg"
                className="flex-1 sm:flex-none"
              >
                Share QR Code
              </Button>
              
              <Button
                onClick={() => {
                  navigator.clipboard.writeText(getRoomJoinUrl())
                  toast.success('Room link copied!')
                }}
                variant="ghost"
                size="lg"
                className="flex-1 sm:flex-none"
              >
                Copy Link
              </Button>
            </div>

            <p className="text-sm text-gray-500">
              Share the Room ID and password with friends to invite them
            </p>
          </Card>
        </motion.div>

        {/* QR Code Modal */}
        <Modal
          isOpen={showQR}
          onClose={() => setShowQR(false)}
          title="Share Room"
          size="md"
        >
          <div className="text-center">
            <div className="bg-white p-4 rounded-lg inline-block mb-4">
              <QRCodeSVG
                value={getRoomJoinUrl()}
                size={200}
                level="M"
                includeMargin={true}
              />
            </div>
            <p className="text-gray-600 mb-4">
              Scan this QR code to join the room
            </p>
            <Button
              onClick={() => {
                navigator.clipboard.writeText(getRoomJoinUrl())
                toast.success('Room link copied!')
              }}
              variant="outline"
              className="w-full"
            >
              Copy Room Link
            </Button>
          </div>
        </Modal>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 py-12 px-4">
      <div className="max-w-2xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          <Card className="p-8">
            {/* Header */}
            <div className="text-center mb-8">
              <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Music className="w-8 h-8 text-primary-600" />
              </div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">
                Create Music Room
              </h1>
              <p className="text-gray-600">
                Set up your collaborative music experience
              </p>
            </div>

            {/* Form */}
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <Input
                label="Room Name"
                {...register('roomName', { 
                  required: 'Room name is required',
                  minLength: { value: 3, message: 'Room name must be at least 3 characters' }
                })}
                error={errors.roomName?.message}
                placeholder="Enter a fun room name"
              />

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <Input
                    label="Session Duration (minutes)"
                    type="number"
                    {...register('lifeTime', { 
                      required: 'Duration is required',
                      min: { value: 15, message: 'Minimum 15 minutes' },
                      max: { value: 480, message: 'Maximum 8 hours' }
                    })}
                    error={errors.lifeTime?.message}
                    min="15"
                    max="480"
                  />
                </div>

                <div>
                  <Input
                    label="Max Users"
                    type="number"
                    {...register('maxUsers', { 
                      required: 'Max users is required',
                      min: { value: 2, message: 'Minimum 2 users' },
                      max: { value: 100, message: 'Maximum 100 users' }
                    })}
                    error={errors.maxUsers?.message}
                    min="2"
                    max="100"
                  />
                </div>
              </div>

              {/* Room Type */}
              <div>
                <label className="label mb-3">Room Type</label>
                <div className="grid grid-cols-2 gap-4">
                  <motion.label
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                    className={`relative cursor-pointer rounded-lg border-2 p-4 transition-all ${
                      watchIsPublic 
                        ? 'border-primary-500 bg-primary-50' 
                        : 'border-gray-200 bg-white hover:border-gray-300'
                    }`}
                  >
                    <input
                      type="radio"
                      value={true}
                      {...register('isPublic')}
                      className="sr-only"
                    />
                    <div className="flex items-center">
                      <Globe className="w-5 h-5 text-primary-600 mr-3" />
                      <div>
                        <div className="font-medium text-gray-900">Public</div>
                        <div className="text-sm text-gray-500">Anyone can discover</div>
                      </div>
                    </div>
                  </motion.label>

                  <motion.label
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                    className={`relative cursor-pointer rounded-lg border-2 p-4 transition-all ${
                      !watchIsPublic 
                        ? 'border-primary-500 bg-primary-50' 
                        : 'border-gray-200 bg-white hover:border-gray-300'
                    }`}
                  >
                    <input
                      type="radio"
                      value={false}
                      {...register('isPublic')}
                      className="sr-only"
                    />
                    <div className="flex items-center">
                      <Lock className="w-5 h-5 text-primary-600 mr-3" />
                      <div>
                        <div className="font-medium text-gray-900">Private</div>
                        <div className="text-sm text-gray-500">Invite only</div>
                      </div>
                    </div>
                  </motion.label>
                </div>
              </div>

              {/* Hidden field for hostUsername */}
              <input type="hidden" {...register('hostUsername')} />

              <Button
                type="submit"
                loading={isCreating}
                size="lg"
                className="w-full"
              >
                Create Room
              </Button>
            </form>
          </Card>
        </motion.div>
      </div>
    </div>
  )
}

export default CreateRoom