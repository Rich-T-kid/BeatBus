from pydantic import BaseModel
from typing import Optional

class SongRequest(BaseModel):
    song_name: str
    artist_name: str
    album_name: str

class SongResponse(BaseModel):
    presigned_url: Optional[str] = None
