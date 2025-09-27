import { useEffect, useRef, useState } from 'react'
import { io } from 'socket.io-client'
import toast from 'react-hot-toast'

const SOCKET_URL = import.meta.env.VITE_SOCKET_URL || 'http://localhost:8080'

const useSocket = (roomId, username) => {
  const [socket, setSocket] = useState(null)
  const [isConnected, setIsConnected] = useState(false)
  const [roomState, setRoomState] = useState(null)
  const socketRef = useRef(null)

  useEffect(() => {
    if (!roomId || !username) return

    // Create socket connection
    const newSocket = io(SOCKET_URL, {
      transports: ['websocket'],
      query: {
        roomId,
        username,
      },
    })

    socketRef.current = newSocket
    setSocket(newSocket)

    // Connection events
    newSocket.on('connect', () => {
      setIsConnected(true)
      console.log('Connected to room:', roomId)
      toast.success('Connected to room!')
    })

    newSocket.on('disconnect', () => {
      setIsConnected(false)
      console.log('Disconnected from room')
      toast.error('Disconnected from room')
    })

    newSocket.on('connect_error', (error) => {
      console.error('Connection error:', error)
      toast.error('Failed to connect to room')
    })

    // Room state updates
    newSocket.on('room_state_update', (state) => {
      setRoomState(state)
    })

    // Queue updates
    newSocket.on('queue_updated', (queue) => {
      setRoomState(prev => ({ ...prev, queue }))
    })

    // Song changes
    newSocket.on('song_changed', (songData) => {
      setRoomState(prev => ({ ...prev, nowPlaying: songData }))
      toast.success(`Now playing: ${songData.stats.title}`)
    })

    // Vote updates
    newSocket.on('vote_update', (voteData) => {
      setRoomState(prev => ({
        ...prev,
        nowPlaying: {
          ...prev.nowPlaying,
          metadata: {
            ...prev.nowPlaying.metadata,
            likes: voteData.likes,
            dislikes: voteData.dislikes,
          }
        }
      }))
    })

    // User events
    newSocket.on('user_joined', (userData) => {
      toast.success(`${userData.username} joined the room`)
      setRoomState(prev => ({
        ...prev,
        numberOfUsers: prev.numberOfUsers + 1
      }))
    })

    newSocket.on('user_left', (userData) => {
      toast(`${userData.username} left the room`)
      setRoomState(prev => ({
        ...prev,
        numberOfUsers: prev.numberOfUsers - 1
      }))
    })

    // Room events
    newSocket.on('room_deleted', () => {
      toast.error('Room has been closed by the host')
      // Redirect to home or handle room closure
      window.location.href = '/'
    })

    // Error handling
    newSocket.on('error', (error) => {
      console.error('Socket error:', error)
      toast.error(error.message || 'Something went wrong')
    })

    // Cleanup on unmount
    return () => {
      newSocket.disconnect()
      setSocket(null)
      setIsConnected(false)
    }
  }, [roomId, username])

  // Socket action methods
  const socketActions = {
    // Join room
    joinRoom: (roomPassword) => {
      socket?.emit('join_room', { roomId, roomPassword, username })
    },

    // Add song to queue
    addSong: (songData) => {
      socket?.emit('add_song', songData)
    },

    // Vote on current song
    voteSong: (action) => {
      socket?.emit('vote_song', { action, songId: roomState?.nowPlaying?.songId })
    },

    // Skip song (host only)
    skipSong: () => {
      socket?.emit('skip_song')
    },

    // Reorder queue (host only)
    reorderQueue: (newOrder) => {
      socket?.emit('reorder_queue', { newOrder })
    },

    // Remove user (host only)
    removeUser: (userId) => {
      socket?.emit('remove_user', { userId })
    },

    // Update room settings (host only)
    updateRoomSettings: (settings) => {
      socket?.emit('update_room_settings', settings)
    },

    // Send message/chat
    sendMessage: (message) => {
      socket?.emit('chat_message', { message, username })
    },
  }

  return {
    socket,
    isConnected,
    roomState,
    ...socketActions,
  }
}

export default useSocket