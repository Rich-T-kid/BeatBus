from fastapi import APIRouter, HTTPException
from models.song import SongRequest, SongResponse
from services.redis_services import get_song_by_metadata, add_song_to_redis
from services import s3_services
from services.scraper_services import scrape_song

router = APIRouter()

@router.post("/songs/search", response_model=SongResponse)
def search_song(song: SongRequest):
    """Get a song data with the given song name, artist name, and album name"""

    result = get_song_by_metadata(song.song_name, song.artist_name, song.album_name)
    if result is None:
        # Try to scrape the song
        s3_key = scrape_song(song.song_name, song.artist_name, song.album_name)
        if s3_key is None:
            raise HTTPException(status_code=404, detail="Song not found")
        else:
            # Add to Redis cache
            add_song_to_redis(song.song_name, song.artist_name, song.album_name)
            result = s3_key
    

    presigned_url = s3_services.create_presigned_url("beatbus-songs", result, 3600)
    
    print(f"🔗 Presigned URL: {presigned_url}")

    return SongResponse(presigned_url=presigned_url)

