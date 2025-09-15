import grpc
from concurrent import futures

import download_pb2_grpc
import download_pb2

from services.redis_services import get_song_by_metadata, add_song_to_redis
from services import s3_services
from services.scraper_services import scrape_song

class DownloadService(download_pb2_grpc.DownloadServiceServicer):
    def HealthCheck(self,request,context):
        print(f"Health check received from: {request.name}")
        return download_pb2.HelloResponse(message="Server is healthy")

    def DownloadSong(self,request,context):
        song_name = request.song_name
        song_artist = request.artist_name
        song_album = request.album_name

        result = get_song_by_metadata(song_name, song_artist, song_album)
        if result is None:
            # Try to scrape the song
            s3_key = scrape_song(song_name, song_artist, song_album)
            if s3_key is None:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details('Song not found')
                return download_pb2.DownloadResponse(download_url="")
                
            else:
                # Add to Redis cache
                add_song_to_redis(song_name, song_artist, song_album)
                result = s3_key
    

        presigned_url = s3_services.create_presigned_url("beatbus-songs", result, 3600)
        return download_pb2.DownloadResponse(download_url=presigned_url)



# Start gRPC server

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    download_pb2_grpc.add_DownloadServiceServicer_to_server(DownloadService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC server started on port 50051")

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        print("Shutting down gRPC server...")
        server.stop(0)

if __name__ == '__main__':
    serve()