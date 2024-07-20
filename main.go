package main

import (
	"log/slog"
)


func main() {
	server := NewServer(":6379")
	err := server.Start()
	if err != nil {
		slog.Error("Failed to start server.")
	}
}
