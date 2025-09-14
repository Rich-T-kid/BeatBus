package storage

/*
I want this to be where we write logs to S3
But for now thats not needed so just to keep note
*/

type SongQueueItem struct {
	Song struct {
		SongID   string                 `bson:"songId"`
		Stats    map[string]interface{} `bson:"stats"`
		Metadata map[string]interface{} `bson:"metadata"`
	} `bson:"song"`
	AlreadyPlayed bool `bson:"alreadyPlayed"`
	Position      int  `bson:"position"`
}
