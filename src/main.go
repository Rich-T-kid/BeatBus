package main

import (
	"BeatBus/server"
	"BeatBus/storage"
	"fmt"
)

func main() {
	_ = storage.NewMessageQueue()

	_ = storage.NewDocumentStore("beatbus")
	fmt.Println("Databases are both connected correctly")
	_ = server.NewServer().StartServer()
}
