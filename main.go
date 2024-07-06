package main

import (
	"log/slog"
)


func main() {
	server := NewServer(":7777")
	err := server.Start()
	if err != nil {
		slog.Error("Failed to start server.")
	}
}