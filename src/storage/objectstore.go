package storage

import (
	"github.com/google/uuid"
)

/*
I want this to be where we write logs to S3
But for now thats not needed so just to keep note
*/
func randomHash() string {
	return uuid.New().String()
}
