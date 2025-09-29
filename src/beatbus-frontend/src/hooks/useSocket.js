import { useEffect, useState } from 'react'
import { io } from 'socket.io-client'
import { mockRoomState } from '../utils/mockData'
import toast from 'react-hot-toast'

const SOCKET_URL = import.meta.env.VITE_SOCKET_URL || 'http://localhost:8080'

const useSocket = (roomId, username) => {
  const [socket, setSocket] = useState(null)
  const [isConnected, setIsConnected] = useState(false)
  const [roomState, setRoomState] = useState(null)
  const [connectionMethod, setConnectionMethod] = useState('websocket')

  useEffect(() => {
    if (!roomId || !username) return

    // Demo mode - use mock data
    if (window.DEMO_MODE) {
      setTimeout(() => {
        setIsConnected(true)
        setRoomState(mockRoomState)
        toast.success('Connected to demo room!')
      }, 1000)
      
      return () => {
        setIsConnected(false)
        setRoomState(null)
      }
    }

    // Real mode - connect to WebSocket
    const newSocket = io(SOCKET_URL, {
      transports: ['websocket'],
      query: { roomId, username },
    })

    setSocket(newSocket)

    newSocket.on('connect', () => {
      setIsConnected(true)
      toast.success('Connected to room!')
    })

    newSocket.on('disconnect', () => {
      setIsConnected(false)
    })

    newSocket.on('room_state_update', (state) => {
      setRoomState(state)
    })

    return () => {
      newSocket.disconnect()
    }
  }, [roomId, username])

  const socketActions = {
    addSong: (songData) => {
      if (window.DEMO_MODE) {
        const newSong = {
          song: {
            songId: `song${Date.now()}`,
            stats: {
              title: songData.songName,
              artist: songData.artistName,
              album: songData.albumName || 'Unknown',
              duration: 180
            },
            metadata: {
              addedBy: songData.addedBy,
              likes: 0,
              dislikes: 0
            }
          },
          alreadyPlayed: false,
          position: (roomState?.queue?.length || 0) + 1
        }
        setRoomState(prev => ({
          ...prev,
          queue: [...(prev?.queue || []), newSong]
        }))
        toast.success('Song added!')
      } else {
        socket?.emit('add_song', songData)
      }
    },

    voteSong: (action) => {
      if (window.DEMO_MODE) {
        console.log('Demo: vote', action)
      } else {
        socket?.emit('vote_song', { action })
      }
    },

    skipSong: () => {
      if (window.DEMO_MODE) {
        if (roomState?.queue?.length > 0) {
          const [nextSong, ...remaining] = roomState.queue
          setRoomState(prev => ({
            ...prev,
            nowPlaying: nextSong.song,
            queue: remaining
          }))
          toast.success(`Now playing: ${nextSong.song.stats.title}`)
        }
      } else {
        socket?.emit('skip_song')
      }
    },

    refreshRoomState: () => {
      if (window.DEMO_MODE) {
        toast.success('Demo: Refreshed')
      }
    },
  }

  return {
    socket,
    isConnected,
    roomState,
    connectionMethod,
    ...socketActions,
  }
}

export default useSocket