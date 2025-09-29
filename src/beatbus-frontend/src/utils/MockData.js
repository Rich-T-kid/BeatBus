// Mock data for demo mode testing

export const mockUser = {
  username: 'testhost',
  id: 'user123'
}

export const mockRoom = {
  roomId: 'demo123',
  roomPassword: 'demo2024',
  hostId: 'testhost',
  roomName: 'Demo Music Room',
  maxUsers: 50,
  isPublic: true
}

export const mockRoomState = {
  roomId: 'demo123',
  numberOfUsers: 5,
  roomSettings: {
    roomName: 'Demo Music Room',
    hostUsername: 'testhost',
    maxUsers: 50,
    isPublic: true
  },
  nowPlaying: {
    songId: 'song001',
    stats: {
      title: 'Blinding Lights',
      artist: 'The Weeknd',
      album: 'After Hours',
      duration: 200
    },
    metadata: {
      addedBy: 'musiclover',
      likes: 15,
      dislikes: 2
    }
  },
  queue: [
    {
      song: {
        songId: 'song002',
        stats: {
          title: 'Levitating',
          artist: 'Dua Lipa',
          album: 'Future Nostalgia',
          duration: 203
        },
        metadata: {
          addedBy: 'popfan',
          likes: 12,
          dislikes: 1
        }
      },
      alreadyPlayed: false,
      position: 1
    },
    {
      song: {
        songId: 'song003',
        stats: {
          title: 'Good 4 U',
          artist: 'Olivia Rodrigo',
          album: 'SOUR',
          duration: 178
        },
        metadata: {
          addedBy: 'rocklover',
          likes: 8,
          dislikes: 0
        }
      },
      alreadyPlayed: false,
      position: 2
    }
  ]
}

export const mockMetrics = {
  roomsize: 5,
  queueLength: 2,
  mostLikedSongs: [
    {
      songName: 'Blinding Lights',
      artistName: 'The Weeknd',
      likes: 15,
      dislikes: 2
    },
    {
      songName: 'Levitating',
      artistName: 'Dua Lipa',
      likes: 12,
      dislikes: 1
    }
  ],
  mostDislikedSongs: [
    {
      songName: 'Blinding Lights',
      artistName: 'The Weeknd',
      likes: 15,
      dislikes: 2
    }
  ],
  userWithMostLikes: 'testhost',
  userWithMostDislikes: 'critic123'
}

export const mockHistory = [
  {
    songId: 'song000',
    stats: {
      title: 'drivers license',
      artist: 'Olivia Rodrigo',
      duration: 242
    },
    metadata: {
      addedBy: 'testhost',
      likes: 25,
      dislikes: 2
    }
  }
]