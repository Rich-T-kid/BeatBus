import { motion } from 'framer-motion'
import { Music, Clock, User, SkipForward } from 'lucide-react'

const CurrentSong = ({ song, isHost, onSkip }) => {
  if (!song) {
    return (
      <div className="text-center py-12">
        <Music className="w-16 h-16 text-gray-300 mx-auto mb-4" />
        <h3 className="text-xl font-medium text-gray-500 mb-2">
          No song playing
        </h3>
        <p className="text-gray-400">
          Add a song to the queue to get started!
        </p>
      </div>
    )
  }

  const formatDuration = (seconds) => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}:${secs.toString().padStart(2, '0')}`
  }

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      className="text-center"
    >
      {/* Album Art Placeholder */}
      <div className="w-48 h-48 mx-auto mb-6 bg-gradient-to-br from-primary-400 to-secondary-500 rounded-2xl flex items-center justify-center shadow-2xl">
        <Music className="w-24 h-24 text-white" />
      </div>
      
      {/* Song Info */}
      <div className="space-y-2 mb-6">
        <h2 className="text-3xl font-bold text-gray-900">
          {song.stats.title}
        </h2>
        <p className="text-xl text-gray-600">
          {song.stats.artist}
        </p>
        {song.stats.album && (
          <p className="text-lg text-gray-500">
            {song.stats.album}
          </p>
        )}
      </div>
      
      {/* Song Metadata */}
      <div className="flex items-center justify-center space-x-6 text-sm text-gray-500">
        {song.stats.duration && (
          <div className="flex items-center space-x-1">
            <Clock className="w-4 h-4" />
            <span>{formatDuration(song.stats.duration)}</span>
          </div>
        )}
        
        <div className="flex items-center space-x-1">
          <User className="w-4 h-4" />
          <span>Added by {song.metadata.addedBy}</span>
        </div>
      </div>
      
      {/* Music Visualizer */}
      <div className="music-visualizer mt-8">
        {[...Array(5)].map((_, i) => (
          <div
            key={i}
            className="music-bar w-1"
            style={{
              height: `${Math.random() * 20 + 4}px`,
              animationDelay: `${i * 0.1}s`
            }}
          />
        ))}
      </div>
      
      {/* Host Controls */}
      {isHost && (
        <div className="mt-6">
          <button
            onClick={onSkip}
            className="btn-outline"
          >
            <SkipForward className="w-4 h-4 mr-2" />
            Skip Song
          </button>
        </div>
      )}
    </motion.div>
  )
}

export default CurrentSong