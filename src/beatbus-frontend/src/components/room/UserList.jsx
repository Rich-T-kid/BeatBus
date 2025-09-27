import { motion } from 'framer-motion'
import { Crown, UserX } from 'lucide-react'
import Modal from '../ui/Modal'
import Button from '../ui/Button'

const UserList = ({ users, isHost, onClose, onRemoveUser }) => {
  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title="Room Users"
      size="md"
    >
      <div className="space-y-4">
        {users.length === 0 ? (
          <p className="text-center text-gray-500 py-8">
            No users connected
          </p>
        ) : (
          users.map((user, index) => (
            <motion.div
              key={user.id}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: index * 0.1 }}
              className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-primary-100 rounded-full flex items-center justify-center">
                  <span className="text-primary-600 font-medium">
                    {user.username[0]?.toUpperCase()}
                  </span>
                </div>
                <div>
                  <div className="flex items-center space-x-2">
                    <span className="font-medium text-gray-900">
                      {user.username}
                    </span>
                    {user.isHost && (
                      <Crown className="w-4 h-4 text-yellow-500" />
                    )}
                  </div>
                  <p className="text-sm text-gray-500">
                    {user.isHost ? 'Host' : 'Guest'}
                  </p>
                </div>
              </div>
              
              {isHost && !user.isHost && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onRemoveUser(user.id)}
                  className="text-red-600 hover:text-red-700"
                >
                  <UserX className="w-4 h-4" />
                </Button>
              )}
            </motion.div>
          ))
        )}
      </div>
      
      <div className="flex justify-end pt-4">
        <Button onClick={onClose} variant="outline">
          Close
        </Button>
      </div>
    </Modal>
  )
}

export default UserList