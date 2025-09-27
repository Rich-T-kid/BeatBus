import { useState, useEffect, useMemo } from 'react'
import { useParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  BarChart3,
  Heart,
  ThumbsDown,
  Crown,
  Music,
  Users,
  Download,
  ArrowLeft,
  Trophy,
  Clock,
  TrendingUp,
} from 'lucide-react'
import { metricsAPI } from '../services/api'
import Card from '../components/ui/Card'
import Button from '../components/ui/Button'
import LoadingSpinner from '../components/ui/LoadingSpinner'
import toast from 'react-hot-toast'

/** helpers */
const safeArr = (v) => (Array.isArray(v) ? v : [])
const n = (v) => (typeof v === 'number' && !Number.isNaN(v) ? v : 0)
const pct = (x) => `${Math.round(x)}%`
const fmt1 = (x) => (Number.isFinite(x) ? x.toFixed(1) : '0.0')

const RoomMetrics = () => {
  const { roomId } = useParams()
  const [metrics, setMetrics] = useState(null) // RoomMetrics from GET /metrics/{roomID}
  const [history, setHistory] = useState([])   // SongObject[] from GET /metrics/{roomID}/history
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let isMounted = true
    const load = async () => {
      try {
        const [mRes, hRes] = await Promise.all([
          metricsAPI.getRoomMetrics(roomId),
          metricsAPI.getHistory(roomId),
        ])
        if (!isMounted) return
        setMetrics(mRes?.data ?? {})
        setHistory(hRes?.data ?? [])
      } catch (e) {
        console.error(e)
        toast.error('Failed to load metrics')
      } finally {
        if (isMounted) setLoading(false)
      }
    }
    load()
    return () => { isMounted = false }
  }, [roomId])

  const {
    roomsize = 0,
    queueLength = 0,
    userWithMostLikes,
    userWithMostDislikes,
  } = metrics ?? {}

  /** derive everything displayable from history (which has stats + metadata) */
  const derived = useMemo(() => {
    const H = safeArr(history)
    const totalSongs = H.length
    const totalLikes = H.reduce((s, x) => s + n(x?.metadata?.likes), 0)
    const totalDislikes = H.reduce((s, x) => s + n(x?.metadata?.dislikes), 0)
    const totalVotes = totalLikes + totalDislikes
    const positiveRatio = totalVotes > 0 ? (totalLikes / totalVotes) * 100 : 0
    const songsPerUser = roomsize > 0 ? totalSongs / roomsize : 0
    const avgLikesPerSong = totalSongs > 0 ? totalLikes / totalSongs : 0

    // Sorts using history so we can show title/artist
    const topLiked = [...H].sort((a, b) => n(b?.metadata?.likes) - n(a?.metadata?.likes))
    const topDisliked = [...H].sort((a, b) => n(b?.metadata?.dislikes) - n(a?.metadata?.dislikes))

    return {
      totalSongs,
      totalLikes,
      totalDislikes,
      totalVotes,
      positiveRatio,
      songsPerUser,
      avgLikesPerSong,
      topLiked,
      topDisliked,
    }
  }, [history, roomsize])

  const handleExportPlaylist = async (mostLikedOnly = false) => {
    try {
      const msg = mostLikedOnly
        ? 'Enter your email or phone number to receive the most liked songs:'
        : 'Enter your email or phone number to receive all songs:'
      const contact = window.prompt(msg)
      if (!contact) return

      // Spec note: only SMS may work “for now”. We’ll infer method.
      const method = contact.includes('@') ? 'email' : 'notification'
      if (method === 'email') {
        toast('Heads up: email may not be supported yet — SMS works best right now.', { icon: 'ℹ️' })
      }

      await metricsAPI.sendPlaylist(roomId, {
        userIds: [
          {
            userId: 'current-user',
            method,           // 'email' | 'notification' (your backend treats notification as sms)
            means: contact,   // 'user@example.com' | '9087654321'
            includeMostLikedOnly: mostLikedOnly,
          },
        ],
      })
      toast.success('Playlist sent successfully!')
    } catch (e) {
      console.error(e)
      toast.error('Failed to send playlist')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <LoadingSpinner size="lg" />
          <p className="text-gray-600 mt-4">Loading session metrics...</p>
        </div>
      </div>
    )
  }

  if (!metrics || Object.keys(metrics).length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <BarChart3 className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-gray-900 mb-2">No Metrics Available</h2>
          <p className="text-gray-600 mb-6">This room hasn’t generated analytics yet.</p>
          <div className="space-x-4">
            <Link to={`/room/${roomId}`}><Button variant="primary">Return to Room</Button></Link>
            <Link to="/"><Button variant="outline">Go Home</Button></Link>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-8 gap-4"
        >
          <div className="flex items-center space-x-4">
            <Link
              to={`/room/${roomId}`}
              className="inline-flex items-center text-gray-600 hover:text-gray-900 transition-colors group"
            >
              <ArrowLeft className="w-5 h-5 mr-1 group-hover:-translate-x-1 transition-transform" />
              Back to Room
            </Link>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Session Analytics</h1>
              <p className="text-gray-600">
                Room {roomId} • {n(roomsize)} participants • {derived.totalSongs} songs played
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-3">
            <Button onClick={() => handleExportPlaylist(true)} variant="outline" size="sm">
              <Heart className="w-4 h-4 mr-2 text-green-600" />
              Export Top Songs
            </Button>
            <Button onClick={() => handleExportPlaylist(false)} variant="primary" size="sm">
              <Download className="w-4 h-4 mr-2" />
              Export All
            </Button>
          </div>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Main Column */}
          <div className="lg:col-span-3 space-y-6">
            {/* Overview */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              {[
                { icon: Users, label: 'Participants', value: n(roomsize), color: 'text-blue-600', bg: 'bg-blue-50' },
                { icon: Music, label: 'Songs Played', value: derived.totalSongs, color: 'text-purple-600', bg: 'bg-purple-50' },
                { icon: Heart, label: 'Total Likes', value: derived.totalLikes, color: 'text-green-600', bg: 'bg-green-50' },
                { icon: TrendingUp, label: 'Engagement', value: pct(derived.positiveRatio), color: 'text-indigo-600', bg: 'bg-indigo-50' },
              ].map((s, i) => (
                <motion.div key={s.label} initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: i * 0.1 }}>
                  <Card className="text-center p-6">
                    <div className={`w-12 h-12 ${s.bg} rounded-lg flex items-center justify-center mx-auto mb-3`}>
                      <s.icon className={`w-6 h-6 ${s.color}`} />
                    </div>
                    <p className="text-2xl font-bold text-gray-900 mb-1">{s.value}</p>
                    <p className="text-sm text-gray-600">{s.label}</p>
                  </Card>
                </motion.div>
              ))}
            </div>

            {/* Top Liked (from history) */}
            <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.3 }}>
              <Card>
                <div className="flex items-center justify-between mb-6">
                  <h2 className="text-xl font-bold text-gray-900 flex items-center">
                    <Trophy className="w-6 h-6 text-yellow-500 mr-2" />
                    Top Liked Songs
                  </h2>
                  <span className="text-sm text-gray-500">
                    showing {Math.min(8, derived.topLiked.length)} of {derived.topLiked.length}
                  </span>
                </div>

                <div className="space-y-3">
                  {derived.topLiked.length > 0 ? (
                    derived.topLiked.slice(0, 8).map((song, index) => {
                      const key = `${song?.songId ?? song?.stats?.title ?? 'liked'}-${index}`
                      const likes = n(song?.metadata?.likes)
                      const dislikes = n(song?.metadata?.dislikes)
                      return (
                        <motion.div
                          key={key}
                          initial={{ opacity: 0, x: -20 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: index * 0.06 }}
                          className="flex items-center justify-between p-4 bg-gradient-to-r from-green-50 to-yellow-50 rounded-lg border border-green-100"
                        >
                          <div className="flex items-center space-x-4">
                            <div className="w-8 h-8 bg-yellow-100 rounded-full flex items-center justify-center">
                              <span className="text-yellow-600 font-bold text-sm">#{index + 1}</span>
                            </div>
                            <div>
                              <h4 className="font-semibold text-gray-900">
                                {song?.stats?.title || 'Unknown Song'}
                              </h4>
                              <p className="text-sm text-gray-600">
                                by {song?.stats?.artist || 'Unknown Artist'}
                              </p>
                              {song?.metadata?.addedBy && (
                                <p className="text-xs text-gray-500">Added by {song.metadata.addedBy}</p>
                              )}
                            </div>
                          </div>
                          <div className="flex items-center space-x-3">
                            <div className="flex items-center space-x-1 text-green-600 bg-green-100 px-3 py-1 rounded-full">
                              <Heart className="w-4 h-4" />
                              <span className="font-bold">{likes}</span>
                            </div>
                            {dislikes > 0 && (
                              <div className="flex items-center space-x-1 text-red-600 bg-red-100 px-3 py-1 rounded-full">
                                <ThumbsDown className="w-4 h-4" />
                                <span className="font-medium">{dislikes}</span>
                              </div>
                            )}
                          </div>
                        </motion.div>
                      )
                    })
                  ) : (
                    <div className="text-center py-12">
                      <Heart className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                      <h3 className="text-lg font-medium text-gray-900 mb-2">No Liked Songs Yet</h3>
                      <p className="text-gray-500">Songs will appear here when users start voting</p>
                    </div>
                  )}
                </div>
              </Card>
            </motion.div>

            {/* Most Disliked (from history) */}
            <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.4 }}>
              <Card>
                <div className="flex items-center justify-between mb-6">
                  <h2 className="text-xl font-semibold text-gray-900">Most Disliked Songs</h2>
                  <ThumbsDown className="w-6 h-6 text-red-500" />
                </div>

                <div className="space-y-4">
                  {derived.topDisliked.length > 0 ? (
                    derived.topDisliked.slice(0, 5).map((song, index) => {
                      const key = `${song?.songId ?? song?.stats?.title ?? 'disliked'}-${index}`
                      return (
                        <div key={key} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                          <div className="flex items-center space-x-4">
                            <div className="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
                              <span className="text-red-600 font-bold text-sm">#{index + 1}</span>
                            </div>
                            <div>
                              <h4 className="font-medium text-gray-900">{song?.stats?.title || 'Unknown Song'}</h4>
                              <p className="text-sm text-gray-600">
                                {song?.stats?.artist || 'Unknown Artist'}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center space-x-2 text-red-600">
                            <ThumbsDown className="w-4 h-4" />
                            <span className="font-medium">{n(song?.metadata?.dislikes)}</span>
                          </div>
                        </div>
                      )
                    })
                  ) : (
                    <div className="text-center py-8">
                      <ThumbsDown className="w-12 h-12 text-gray-300 mx-auto mb-4" />
                      <p className="text-gray-500">No disliked songs</p>
                    </div>
                  )}
                </div>
              </Card>
            </motion.div>

            {/* History */}
            <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.5 }}>
              <Card>
                <div className="flex items-center justify-between mb-6">
                  <h2 className="text-xl font-bold text-gray-900 flex items-center">
                    <Clock className="w-6 h-6 text-blue-500 mr-2" />
                    Complete Session History
                  </h2>
                  <span className="text-sm text-gray-500">{derived.totalSongs} total songs</span>
                </div>

                <div className="space-y-2 max-h-96 overflow-y-auto">
                  {safeArr(history).length > 0 ? (
                    safeArr(history).map((song, index) => {
                      const key = `${song?.songId ?? song?.stats?.title ?? 'hist'}-${index}`
                      return (
                        <motion.div
                          key={key}
                          initial={{ opacity: 0 }}
                          animate={{ opacity: 1 }}
                          transition={{ delay: index * 0.02 }}
                          className="flex items-center justify-between p-3 hover:bg-gray-50 rounded-lg border-l-4 border-blue-200"
                        >
                          <div className="flex items-center space-x-3">
                            <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                              <span className="text-blue-600 font-medium text-xs">
                                {derived.totalSongs - index}
                              </span>
                            </div>
                            <div>
                              <h4 className="font-medium text-gray-900 text-sm">
                                {song?.stats?.title || 'Unknown Song'}
                              </h4>
                              <p className="text-xs text-gray-600">
                                {song?.stats?.artist || 'Unknown Artist'}
                                {song?.metadata?.addedBy ? ` • Added by ${song.metadata.addedBy}` : ''}
                              </p>
                              {song?.stats?.album && (
                                <p className="text-xs text-gray-500">Album: {song.stats.album}</p>
                              )}
                            </div>
                          </div>
                          <div className="flex items-center space-x-2">
                            {n(song?.metadata?.likes) > 0 && (
                              <div className="flex items-center space-x-1 text-green-600 text-xs">
                                <Heart className="w-3 h-3" />
                                <span>{n(song?.metadata?.likes)}</span>
                              </div>
                            )}
                            {n(song?.metadata?.dislikes) > 0 && (
                              <div className="flex items-center space-x-1 text-red-600 text-xs">
                                <ThumbsDown className="w-3 h-3" />
                                <span>{n(song?.metadata?.dislikes)}</span>
                              </div>
                            )}
                          </div>
                        </motion.div>
                      )
                    })
                  ) : (
                    <div className="text-center py-12">
                      <Clock className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                      <h3 className="text-lg font-medium text-gray-900 mb-2">No History Yet</h3>
                      <p className="text-gray-500">Song history will appear as music is played</p>
                    </div>
                  )}
                </div>
              </Card>
            </motion.div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Hall of Fame */}
            <motion.div initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} transition={{ delay: 0.6 }}>
              <Card>
                <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center">
                  <Crown className="w-5 h-5 text-yellow-500 mr-2" />
                  Hall of Fame
                </h3>

                <div className="space-y-4">
                  {userWithMostLikes ? (
                    <div className="p-4 bg-gradient-to-r from-yellow-50 to-orange-50 rounded-lg border border-yellow-200">
                      <div className="text-center">
                        <Crown className="w-8 h-8 text-yellow-500 mx-auto mb-2" />
                        <h4 className="font-bold text-gray-900">Mr. Put 'Em On</h4>
                        <p className="text-lg font-semibold text-yellow-700">{userWithMostLikes}</p>
                        <p className="text-xs text-gray-600 mt-1">Most liked song contributions</p>
                      </div>
                    </div>
                  ) : (
                    <div className="p-4 bg-gray-50 rounded-lg text-center">
                      <Crown className="w-8 h-8 text-gray-300 mx-auto mb-2" />
                      <p className="text-sm text-gray-500">No Mr. Put 'Em On yet</p>
                    </div>
                  )}

                  {userWithMostDislikes && (
                    <div className="p-3 bg-gray-50 rounded-lg border border-gray-200">
                      <div className="flex items-center space-x-2">
                        <ThumbsDown className="w-5 h-5 text-gray-500" />
                        <div>
                          <p className="font-medium text-gray-900 text-sm">Most Discerning</p>
                          <p className="text-sm text-gray-700">{userWithMostDislikes}</p>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              </Card>
            </motion.div>

            {/* Quick Stats */}
            <motion.div initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} transition={{ delay: 0.7 }}>
              <Card>
                <h3 className="text-lg font-bold text-gray-900 mb-4">Quick Stats</h3>
                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600 text-sm">Songs in Queue</span>
                    <span className="font-semibold">{n(queueLength)}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600 text-sm">Total Votes</span>
                    <span className="font-semibold">{derived.totalVotes}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600 text-sm">Approval Rate</span>
                    <span className="font-semibold text-green-600">{pct(derived.positiveRatio)}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600 text-sm">Songs per User</span>
                    <span className="font-semibold">{fmt1(derived.songsPerUser)}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600 text-sm">Avg. Likes per Song</span>
                    <span className="font-semibold">{fmt1(derived.avgLikesPerSong)}</span>
                  </div>
                </div>
              </Card>
            </motion.div>

            {/* Export */}
            <motion.div initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} transition={{ delay: 0.8 }}>
              <Card>
                <h3 className="text-lg font-bold text-gray-900 mb-4">Export Playlist</h3>
                <div className="space-y-3">
                  <p className="text-sm text-gray-600">Save and share the music from this session</p>
                  <Button onClick={() => handleExportPlaylist(true)} variant="outline" className="w-full justify-start" size="sm">
                    <Heart className="w-4 h-4 mr-2 text-green-600" />
                    Top Liked Songs ({Math.min(8, derived.topLiked.length)})
                  </Button>
                  <Button onClick={() => handleExportPlaylist(false)} variant="outline" className="w-full justify-start" size="sm">
                    <Music className="w-4 h-4 mr-2" />
                    All Songs ({derived.totalSongs})
                  </Button>
                  <p className="text-xs text-gray-500">Enter your email or phone to receive the playlist</p>
                </div>
              </Card>
            </motion.div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default RoomMetrics
