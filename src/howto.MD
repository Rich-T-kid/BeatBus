# BeatBus API Guide (Frontend Reference)

This document explains how the frontend should use the Swagger documentation to interact with the BeatBus API. It covers the lifecycle of a room, user interactions, and session metrics.

# note to self: 
## make sure to increment song count and to set the position of each song object to the song count + 1


---

## 1. Host Creates a Room

**Endpoint:**  
`POST /rooms`

**Description:**  
- Creates a new logical room.  
- Returns:
  - `roomId` (unique room identifier)  
  - `roomPassword` (used for joining, also embedded in QR code)  
  - `accessToken` (JWT with host privileges)  
- Generates a QR code for sharing with friends.  
- Host can configure rules:
  - Max users  
  - Whether the room is public/private  
  - Song playback constraints (e.g., number of songs per user in a row)  

**Frontend Notes:**  
- Store the host’s JWT securely; only the host can call host-tagged endpoints.  
- Display the QR code to users for joining.  

---

## 2. Users Join a Room

**Endpoint:**  
`GET /rooms/{roomId}`  

**Parameters:**  
- `roomPassword` (query param, embedded in QR code)  
- `username` (query param, required)  

**Description:**  
- Validates the room password.  
- Users must provide at least a name.  
- Users may also provide a phone number if they want to receive the playlist at the end.  

**Frontend Notes:**  
- If the room is private, the QR code will already contain the room password.  
- Ensure name is captured before making the join request.  

---

## 3. Host Privileges

- Only the host (with the host JWT) can:  
  - Update room settings (`PUT /rooms`)  
  - Delete/end the room (`DELETE /rooms`)  
  - Reorder the queue (`PUT /queues/{roomId}/playlist`)  
  - Push out playlists (`POST /metrics/{roomId}/playlist/send`)  

**Endpoint Example:**  
`GET /metrics/{roomId}` – retrieves live room metrics for the host.  

**Frontend Notes:**  
- All endpoints tagged with `Host` require the host’s JWT.  
- Regular users calling these endpoints should expect `401 Unauthorized`.  

---

## 4. Queue Management and Interactions

**Endpoints:**  
- `POST /queues/{roomId}/playlist` – Add a song to the queue.  
- `GET /queues/{roomId}/playlist` – View current queue.  
- `PUT /queues/{roomId}/playlist` – Reorder queue (host only).  

**Rules:**  
- Room rules define how many songs a user can queue.  
- Users can only skip their own songs.  
- Host can skip or reorder any song.  
- During playback:
  - Users can like/dislike songs (`POST /metrics/{roomId}` with `action: like` or `dislike`).  
  - Likes/dislikes update session metrics.  

**Frontend Notes:**  
- Show users only the actions they’re allowed (skip own song, like/dislike any).  
- Queue order changes should be strictly managed on the frontend before sending updates to the backend.  

---

## 5. Ending a Session

**Endpoint:**  
`DELETE /rooms`  

**Description:**  
- Ends the session and kicks users out.  
- Returns final metrics, including **Mr. Put On** (the user with the most liked songs).  
- Host can push the playlist to users (`POST /metrics/{roomId}/playlist/send`), choosing:
  - All songs  
  - Only most liked songs  

**Frontend Notes:**  
- Display Mr. Put On to all participants when the session ends.  
- Provide an option for users to supply contact details (e.g., phone number) if they want playlist stats.  

---

## 6. Real-Time Updates

**Endpoint:**  
`GET /rooms/{roomId}/state`  

**Description:**  
- Returns current room state:
  - Now playing  
  - Queue  
  - Number of users  
  - Room settings  
- Designed to be implemented with **WebSockets** for continuous updates.  

**Frontend Notes:**  
- Initially, you can poll this endpoint for testing.  
- In production, establish a WebSocket connection to keep state synced live.  

---

## 7. Metrics and History

**Endpoints:**  
- `GET /metrics/{roomId}` – View current session metrics.  
- `GET /metrics/{roomId}/history` – View history of played songs (with likes/dislikes).  
- `POST /metrics/{roomId}` – Send like/dislike.  

**Frontend Notes:**  
- Show real-time updates of likes/dislikes as users interact.  
- Metrics are aggregated server-side, so frontend should always refresh state from the API rather than maintaining separate counters.  

---

# Key Takeaways

- **Host = JWT with elevated permissions**. Only host can call `Host` endpoints.  
- **QR codes** embed `roomPassword` for easy joining.  
- **Users** provide at least a username to join; optional phone/email for playlist delivery.  
- **Queue rules** are enforced server-side; frontend must respect them to avoid failed requests.  
- **Real-time updates** will ultimately use WebSockets (`/rooms/{roomId}/state`).  
