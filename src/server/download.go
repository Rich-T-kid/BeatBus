package server

import (
	pb "BeatBus/internal/grpc"
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DownloadQueue struct{}

func NewDownloadQueue() *DownloadQueue {
	return &DownloadQueue{}
}

// Send GRPC call to download song from youtube to S3 bucket
func (dq *DownloadQueue) RetrieveSong(s AddSongRequest) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.DownloadServerIP, cfg.DownloadServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()

	c := pb.NewDownloadServiceClient(conn)
	resp, _ := c.DownloadSong(context.Background(), &pb.DownloadRequest{
		SongName:   s.SongName,
		ArtistName: s.ArtistName,
		AlbumName:  s.AlbumName,
	})
	fmt.Printf("Download response: %v\n", resp)

}

func SendToS3(content io.Reader) error {
	fmt.Println("sending to s3...")
	return nil
}
