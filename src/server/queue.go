package server

import "fmt"

type queueManager struct{}

type DownloadQueue struct{}

func NewDownloadQueue() *DownloadQueue {
	return &DownloadQueue{}
}

// Send GRPC call to download song from youtube to S3 bucket
func (dq *DownloadQueue) RetrieveSong(s AddSongRequest) error {
	fmt.Printf("simulating downloading song: %s by %s\n", s.SongName, s.ArtistName)
	return nil
}
