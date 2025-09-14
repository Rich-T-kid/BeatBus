from redis import Redis
from config import config

r = Redis(host=config.REDIS_HOST, port=config.REDIS_PORT, db=config.REDIS_DB, decode_responses=True)

def get_song_by_metadata(song_name, artist_name, album_name):
    redis_key = f"songs/{song_name.lower().replace(' ', '_')}_{artist_name.lower().replace(' ', '_')}_{album_name.lower().replace(' ', '_')}.mp3"
    
    if r.exists(redis_key):
        return redis_key
    else:
        return None

def add_song_to_redis(song_name, artist_name, album_name):
    redis_key = f"songs/{song_name.lower().replace(' ', '_')}_{artist_name.lower().replace(' ', '_')}_{album_name.lower().replace(' ', '_')}.mp3"
    r.set(redis_key, "exists")