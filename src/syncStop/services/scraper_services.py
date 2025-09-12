import yt_dlp
import os
from config import config
from .s3_services import upload_file_to_s3

# added

def scrape_song(song_name, artist_name, album_name):
    ydl_opts = {
        'format': 'bestaudio/best',
        'outtmpl': 'downloaded_song.%(ext)s',  # Save in current directory
        'restrictfilenames': True,  # Remove spaces and special characters
        'nocheckcertificate': True,  # Bypass SSL certificate verification
        'postprocessors': [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': 'mp3',
            'preferredquality': '192',
        }],
    }
    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        ydl.download([f"ytsearch1:{song_name} {artist_name} {album_name}"])
        
        # Simple filename - will be downloaded_song.mp3
        mp3_path = "downloaded_song.mp3"
        
        s3_key = f"songs/{song_name.lower().replace(' ', '_')}_{artist_name.lower().replace(' ', '_')}_{album_name.lower().replace(' ', '_')}.mp3"
        upload_file_to_s3("beatbus-songs", mp3_path, s3_key)
        os.remove(mp3_path)  # Clean up local file
        return s3_key
    
    