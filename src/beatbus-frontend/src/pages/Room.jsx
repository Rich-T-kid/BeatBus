import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  Music, 
  Users, 
  Heart, 
  ThumbsDown, 
  SkipForward, 
  Plus,
  Settings,
  Share2,
  BarChart3,
  Crown
} from 'lucide-react'
import useSocket from '../hooks/useSocket'
import useAuthStore from '../store/authStore'
import CurrentSong from '../components/room/CurrentSong'
import QueueList from '../components/room/QueueList'
import VoteButtons from '../components/room/VoteButtons'
import UserList from '../components/room/UserList'
import AddSongModal from '../components/room/AddSongModal'
import Card from '../components/ui/Card'
import Button from '../components/ui/Button'
import toast from 'react-hot-toast'

const Room = () => {
  const { roomId } = useParams()
  const navigate = useNavigate()
  const { user, isAuthenticated } = useAuthStore()
  
  // Get username from localStorage for guests or from auth for hosts
  const username = isAuthenticated ? user?.username : localStorage.getItem('guestUsername')
  
  const {
    isConnected,
    roomState,
    addSong,
    voteSong,
    skipSong,
    reorderQueue,
    removeUser,
    updateRoomSettings
  } = useSocket(roomId, username)
  
  const [showAddSong, setShowAddSong] = useState(false)
  const [showSettings, setShowSettings] = useState(false)
  const [showUsers, setShowUsers] = useState(false)
  
  // Check if current user is the host
  const isHost = isAuthenticated && user?.username === roomState?.roomSettings?.hostUsername
  
  useEffect(() => {
    if (!username) {
      toast.error('Please join the room first')
      navigate('/join')
      return
    }
    
    if (!isConnected && roomState === null) {
      toast.loading('Connecting to room...', { id: 'connecting' })
    } else if (isConnected) {
      toast.dismiss('connecting')
    }
  }, [isConnected, roomState, username, navigate])

  const handleVote = (action) => {
    if (!roomState?.nowPlaying) {
      toast.error('No song is currently playing')
      return
    }
    voteSong(action)
  }

  const handleSkip = () => {
    if (!isHost) {
      toast.error('Only the host can skip songs')
      return
    }
    skipSong()
    toast.success('Song skipped')
  }

  const handleAddSong = (songData) => {
    addSong({
      ...songData,
      addedBy: username
    })
    setShowAddSong(false)
    toast.success('Song added to queue!')
  }

  if (!isConnected && !roomState) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="spinner w-12 h-12 mx-auto mb-4"></div>
          <p className="text-lg text-gray-600">Connecting to room...</p>
        </div>
      </div>
    )
  }

  if (!roomState) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-lg text-gray-600">Room not found</p>
          <Button 
            onClick={() => navigate('/')}
            className="mt-4"
          >
            Go Home
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Room Header */}
      <div className="room-header shadow-lg">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2">
                <Music className="w-8 h-8" />
                <div>
                  <h1 className="text-2xl font-bold">
                    {roomState.roomSettings?.roomName || 'Music Room'}
                  </h1>
                  <p className="text-white/80">
                    Room ID: {roomId}
                  </p>
                </div>
              </div>
              
              {isHost && (
                <div className="flex items-center space-x-1 bg-yellow-500/20 px-3 py-1 rounded-full">
                  <Crown className="w-4 h-4" />
                  <span className="text-sm font-medium">Host</span>
                </div>
              )}
            </div>
            
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2 bg-white/10 px-3 py-2 rounded-lg">
                <Users className="w-5 h-5" />
                <span className="font-medium">{roomState.numberOfUsers}</span>
              </div>
              
              <div className="flex items-center space-x-2">
                <Button
                  onClick={() => setShowUsers(true)}
                  variant="ghost"
                  className="text-white hover:bg-white/10"
                >
                  <Users className="w-5 h-5" />
                </Button>
                
                <Button
                  onClick={() => navigate(`/metrics/${roomId}`)}
                  variant="ghost"
                  className="text-white hover:bg-white/10"
                >
                  <BarChart3 className="w-5 h-5" />
                </Button>
                
                <Button
                  onClick={() => {
                    navigator.share?.({
                      title: 'Join my BeatBus room!',
                      text: `Join "${roomState.roomSettings?.roomName}" on BeatBus`,
                      url: window.location.href
                    }) || navigator.clipboard.writeText(window.location.href).then(() => {
                      toast.success('Room link copied!')
                    })
                  }}
                  variant="ghost"
                  className="text-white hover:bg-white/10"
                >
                  <Share2 className="w-5 h-5" />
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Current Song & Controls */}
          <div className="lg:col-span-2 space-y-6">
            {/* Now Playing */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="card"
            >
              <CurrentSong 
                song={roomState.nowPlaying}
                isHost={isHost}
                onSkip={handleSkip}
              />
              
              <div className="mt-6 flex items-center justify-between">
                <VoteButtons
                  song={roomState.nowPlaying}
                  onVote={handleVote}
                  username={username}
                />
                
                <div className="flex items-center space-x-4">
                  <Button
                    onClick={() => setShowAddSong(true)}
                    variant="primary"
                  >
                    <Plus className="w-5 h-5 mr-2" />
                    Add Song
                  </Button>
                  
                  {isHost && (
                    <Button
                      onClick={handleSkip}
                      variant="outline"
                    >
                      <SkipForward className="w-5 h-5 mr-2" />
                      Skip
                    </Button>
                  )}
                </div>
              </div>
            </motion.div>
            
            {/* Queue */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
              className="card"
            >
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-xl font-semibold text-gray-900">
                  Queue ({roomState.queue?.length || 0})
                </h2>
                {isHost && roomState.queue?.length > 0 && (
                  <Button variant="ghost" size="sm">
                    Shuffle Queue
                  </Button>
                )}
              </div>
              
              <QueueList
                queue={roomState.queue || []}
                isHost={isHost}
                onReorder={reorderQueue}
                currentUser={username}
              />
            </motion.div>
          </div>
          
          {/* Sidebar */}
          <div className="space-y-6">
            {/* Room Stats */}
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.2 }}
              className="card"
            >
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                Room Stats
              </h3>
              
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Active Users</span>
                  <span className="font-semibold">{roomState.numberOfUsers}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Songs in Queue</span>
                  <span className="font-semibold">{roomState.queue?.length || 0}</span>
                </div>
                {roomState.nowPlaying && (
                  <>
                    <div className="flex items-center justify-between">
                      <span className="text-gray-600">Current Likes</span>
                      <span className="flex items-center text-green-600">
                        <Heart className="w-4 h-4 mr-1" />
                        {roomState.nowPlaying.metadata?.likes || 0}
                      </span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-gray-600">Current Dislikes</span>
                      <span className="flex items-center text-red-600">
                        <ThumbsDown className="w-4 h-4 mr-1" />
                        {roomState.nowPlaying.metadata?.dislikes || 0}
                      </span>
                    </div>
                  </>
                )}
              </div>
            </motion.div>
            
            {/* Online Users Preview */}
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.3 }}
              className="card"
            >
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-gray-900">
                  Online Users
                </h3>
                <Button
                  onClick={() => setShowUsers(true)}
                  variant="ghost"
                  size="sm"
                  className="text-primary-600 hover:text-primary-700"
                >
                  View All
                </Button>
              </div>
              
              <div className="space-y-2">
                <div className="flex items-center space-x-3">
                  <div className="w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center">
                    <span className="text-primary-600 font-medium text-sm">
                      {username?.[0]?.toUpperCase()}
                    </span>
                  </div>
                  <span className="text-gray-900 font-medium">
                    {username} {isHost && '(Host)'}
                  </span>
                </div>
                
                {/* Show other users when available */}
                <div className="text-sm text-gray-500 pt-2">
                  +{Math.max(0, (roomState.numberOfUsers || 1) - 1)} other{roomState.numberOfUsers > 2 ? 's' : ''}
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </div>
      
      {/* Modals */}
      <AnimatePresence>
        {showAddSong && (
          <AddSongModal
            onClose={() => setShowAddSong(false)}
            onAddSong={handleAddSong}
          />
        )}
        
        {showUsers && (
          <UserList
            users={[]} // Populate with actual user list
            isHost={isHost}
            onClose={() => setShowUsers(false)}
            onRemoveUser={removeUser}
          />
        )}
      </AnimatePresence>
    </div>
  )
}

export default Room