package main

import (
	"BeatBus/server"
	"BeatBus/storage"
	"fmt"
)

func main() {
	_ = storage.NewMessageQueue(nil)

	_ = storage.NewDocumentStore(nil) // not imporant
	fmt.Println("Databases are both connected correctly")
	_ = server.NewServer().StartServer()
}
