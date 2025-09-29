import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Music, Users, Heart, QrCode, Play, Sparkles } from 'lucide-react'
import useAuthStore from '../store/authStore'

const Home = () => {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuthStore()

  const features = [
    {
      icon: <Users className="w-8 h-8" />,
      title: "Real-time Collaboration",
      description: "Create rooms and invite friends to build playlists together",
      color: "blue"
    },
    {
      icon: <Heart className="w-8 h-8" />,
      title: "Vote on Music",
      description: "Like and dislike songs to influence what plays next",
      color: "green"
    },
    {
      icon: <QrCode className="w-8 h-8" />,
      title: "Easy Joining",
      description: "Join rooms instantly with QR codes - no signup required",
      color: "blue"
    },
    {
      icon: <Sparkles className="w-8 h-8" />,
      title: "Session Analytics",
      description: "Get insights and export your favorite tracks",
      color: "green"
    }
  ]

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative overflow-hidden bg-gradient-to-br from-slate-800 via-slate-900 to-slate-950">
        <div className="absolute inset-0 bg-black/20"></div>
        
        {/* Animated background elements */}
        <div className="absolute inset-0">
          <div className="absolute top-1/4 left-1/4 w-64 h-64 bg-blue-500/10 rounded-full blur-3xl animate-pulse" style={{ animationDuration: '4s' }}></div>
          <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-green-500/10 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '2s', animationDuration: '4s' }}></div>
        </div>
        
        <div className="relative max-w-7xl mx-auto px-4 py-24 sm:px-6 lg:px-8">
          <motion.div 
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-center"
          >
            {/* Logo */}
            <motion.div 
              className="flex items-center justify-center mb-8"
              whileHover={{ scale: 1.05 }}
              transition={{ type: "spring", stiffness: 300 }}
            >
              <div className="relative">
                <Music className="w-16 h-16 text-white" />
                <div className="absolute -top-2 -right-2 w-6 h-6 bg-green-400 rounded-full animate-bounce"></div>
              </div>
              <h1 className="ml-4 text-5xl sm:text-6xl font-bold text-white">
                BeatBus
              </h1>
            </motion.div>
            
            <motion.p 
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.3, duration: 0.8 }}
              className="text-xl sm:text-2xl text-white/90 font-medium mb-4"
            >
              Share and Experience Music Together
            </motion.p>
            
            <motion.p 
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.5, duration: 0.8 }}
              className="text-lg text-white/80 mb-12 max-w-2xl mx-auto"
            >
              Create collaborative music rooms, vote on tracks, and discover new music with friends or strangers in real-time.
            </motion.p>
            
            {/* CTA Buttons */}
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.7, duration: 0.8 }}
              className="flex flex-col sm:flex-row gap-4 justify-center items-center"
            >
              {isAuthenticated ? (
                <Link
                  to="/create"
                  className="inline-flex items-center bg-white text-slate-900 hover:bg-gray-100 text-lg px-8 py-4 rounded-xl font-semibold transform hover:scale-105 transition-all duration-200 shadow-lg hover:shadow-xl"
                >
                  <Play className="w-5 h-5 mr-2" />
                  Create Room
                </Link>
              ) : (
                <Link
                  to="/login"
                  className="inline-flex items-center bg-white text-slate-900 hover:bg-gray-100 text-lg px-8 py-4 rounded-xl font-semibold transform hover:scale-105 transition-all duration-200 shadow-lg hover:shadow-xl"
                >
                  <Play className="w-5 h-5 mr-2" />
                  Get Started
                </Link>
              )}
              
              <Link
                to="/join"
                className="inline-flex items-center border-2 border-white/30 bg-white/10 backdrop-blur-sm text-white hover:bg-white hover:text-slate-900 text-lg px-8 py-4 rounded-xl font-semibold transform hover:scale-105 transition-all duration-200"
              >
                <QrCode className="w-5 h-5 mr-2" />
                Join Room
              </Link>
            </motion.div>
          </motion.div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-24 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div 
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-4xl font-bold text-gray-900 mb-4">
              Why Choose BeatBus?
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              Experience music like never before with real-time collaboration and social voting
            </p>
          </motion.div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 30 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1, duration: 0.8 }}
                viewport={{ once: true }}
                whileHover={{ y: -5 }}
                className="bg-white rounded-xl border border-gray-200 p-6 shadow-sm text-center hover:shadow-lg transition-all duration-200"
              >
                <div className={`inline-flex items-center justify-center w-16 h-16 ${
                  feature.color === 'blue' ? 'bg-blue-50 text-blue-600' : 'bg-green-50 text-green-600'
                } rounded-xl mb-4`}>
                  {feature.icon}
                </div>
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  {feature.title}
                </h3>
                <p className="text-gray-600">
                  {feature.description}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* How it Works Section */}
      <section className="py-24 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div 
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-4xl font-bold text-gray-900 mb-4">
              How It Works
            </h2>
            <p className="text-xl text-gray-600">
              Get started in three simple steps
            </p>
          </motion.div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
            {[
              {
                step: "01",
                title: "Create or Join",
                description: "Host creates a room or guests join with a QR code",
                color: "blue"
              },
              {
                step: "02", 
                title: "Add Music",
                description: "Everyone adds their favorite songs to the shared queue",
                color: "green"
              },
              {
                step: "03",
                title: "Vote & Enjoy",
                description: "Vote on tracks and enjoy music together in real-time",
                color: "blue"
              }
            ].map((item, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, x: index % 2 === 0 ? -30 : 30 }}
                whileInView={{ opacity: 1, x: 0 }}
                transition={{ delay: index * 0.2, duration: 0.8 }}
                viewport={{ once: true }}
                className="text-center"
              >
                <div className={`inline-flex items-center justify-center w-16 h-16 ${
                  item.color === 'blue' ? 'bg-blue-100 text-blue-600' : 'bg-green-100 text-green-600'
                } rounded-full text-2xl font-bold mb-4`}>
                  {item.step}
                </div>
                <h3 className="text-2xl font-semibold text-gray-900 mb-4">
                  {item.title}
                </h3>
                <p className="text-gray-600 text-lg">
                  {item.description}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-24 bg-gradient-to-r from-blue-600 to-green-600">
        <div className="max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
          >
            <h2 className="text-4xl font-bold text-white mb-4">
              Ready to Start Your Musical Journey?
            </h2>
            <p className="text-xl text-white/90 mb-8">
              Join thousands of music lovers creating unforgettable shared experiences
            </p>
            
            <motion.div
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <Link
                to={isAuthenticated ? "/create" : "/login"}
                className="inline-flex items-center bg-white text-slate-900 hover:bg-gray-100 px-8 py-4 rounded-xl text-lg font-semibold transition-all duration-200 shadow-lg"
              >
                Get Started Now
                <Play className="w-5 h-5 ml-2" />
              </Link>
            </motion.div>
          </motion.div>
        </div>
      </section>
    </div>
  )
}

export default Home