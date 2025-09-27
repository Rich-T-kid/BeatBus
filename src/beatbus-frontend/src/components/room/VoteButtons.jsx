import { useState, useEffect } from 'react'
import { Heart, ThumbsDown } from 'lucide-react'
import { motion } from 'framer-motion'

const VoteButtons = ({ song, onVote, username }) => {
  const [userVote, setUserVote] = useState(null) // 'like', 'dislike', or null
  
  useEffect(() => {
    // Reset vote when song changes
    setUserVote(null)
  }, [song?.songId])
  
  const handleVote = (action) => {
    if (!song) return
    
    let newAction = action
    
    // If user already voted the same way, remove the vote
    if (userVote === action) {
      newAction = `remove-${action}`
      setUserVote(null)
    } else {
      setUserVote(action)
    }
    
    onVote(newAction)
  }
  
  if (!song) return null
  
  return (
    <div className="flex items-center space-x-4">
      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={() => handleVote('like')}
        className={`vote-btn ${userVote === 'like' ? 'liked' : ''}`}
      >
        <Heart 
          className={`w-5 h-5 ${userVote === 'like' ? 'fill-current' : ''}`} 
        />
        <span>{song.metadata.likes || 0}</span>
      </motion.button>
      
      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={() => handleVote('dislike')}
        className={`vote-btn ${userVote === 'dislike' ? 'disliked' : ''}`}
      >
        <ThumbsDown 
          className={`w-5 h-5 ${userVote === 'dislike' ? 'fill-current' : ''}`} 
        />
        <span>{song.metadata.dislikes || 0}</span>
      </motion.button>
    </div>
  )
}

export default VoteButtons