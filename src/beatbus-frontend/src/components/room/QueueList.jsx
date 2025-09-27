import { motion, Reorder } from 'framer-motion'
import { Music, GripVertical, X, Play, Heart } from 'lucide-react'
import { useState, useEffect } from 'react'

const QueueList = ({ queue, isHost, onReorder, currentUser }) => {
  const [items, setItems] = useState(queue)
  
  useEffect(() => {
    setItems(queue)
  }, [queue])
  
  const handleReorder = (newOrder) => {
    setItems(newOrder)
    if (isHost) {
      const songIds = newOrder.map(item => item.song.songId)
      onReorder(songIds)
    }
  }
  
  const canRemoveSong = (song) => {
    return isHost || song.metadata.addedBy === currentUser
  }
  
  if (!queue || queue.length === 0) {
    return (
      <div className="text-center py-12">
        <Music className="w-12 h-12 text-gray-300 mx-auto mb-4" />
        <p className="text-gray-500">
          Queue is empty. Add some songs to get started!
        </p>
      </div>
    )
  }
  
  const QueueItem = ({ item, index }) => (
    <motion.div
      layout
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className={`queue-item ${item.alreadyPlayed ? 'opacity-50' : ''}`}
    >
      {/* Drag Handle (Host only) */}
      {isHost && !item.alreadyPlayed && (
        <div className="cursor-grab active:cursor-grabbing">
          <GripVertical className="w-5 h-5 text-gray-400" />
        </div>
      )}
      
      {/* Position */}
      <div className="flex-shrink-0 w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center text-sm font-medium text-gray-600">
        {item.alreadyPlayed ? (
          <Play className="w-4 h-4" />
        ) : (
          index + 1
        )}
      </div>
      
      {/* Song Info */}
      <div className="flex-1 min-w-0">
        <h4 className="font-medium text-gray-900 truncate">
          {item.song.stats.title}
        </h4>
        <p className="text-sm text-gray-600 truncate">
          {item.song.stats.artist}
        </p>
        <p className="text-xs text-gray-500">
          Added by {item.song.metadata.addedBy}
        </p>
      </div>
      
      {/* Actions */}
      <div className="flex items-center space-x-2">
        {item.song.metadata.likes > 0 && (
          <div className="flex items-center space-x-1 text-green-600">
            <Heart className="w-4 h-4" />
            <span className="text-sm">{item.song.metadata.likes}</span>
          </div>
        )}
        
        {canRemoveSong(item.song) && !item.alreadyPlayed && (
          <button
            onClick={() => {
              // Handle song removal
              console.log('Remove song:', item.song.songId)
            }}
            className="text-red-500 hover:text-red-700 p-1"
          >
            <X className="w-4 h-4" />
          </button>
        )}
      </div>
    </motion.div>
  )
  
  if (isHost) {
    return (
      <Reorder.Group
        axis="y"
        values={items}
        onReorder={handleReorder}
        className="space-y-3"
      >
        {items.map((item, index) => (
          <Reorder.Item key={item.song.songId} value={item}>
            <QueueItem item={item} index={index} />
          </Reorder.Item>
        ))}
      </Reorder.Group>
    )
  }
  
  return (
    <div className="space-y-3">
      {items.map((item, index) => (
        <QueueItem key={item.song.songId} item={item} index={index} />
      ))}
    </div>
  )
}

export default QueueList