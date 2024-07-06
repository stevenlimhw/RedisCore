package main

import (
	"fmt"
	"log/slog"
	"net"
)

type Peer struct {
	conn net.Conn
}

func (peer *Peer) readMessages() {
	slog.Info("Reading messages sent to peer...")
	for {
		buffer := make([]byte, 1024)
		size, err := peer.conn.Read(buffer)
		if err != nil {
			slog.Error("Error in reading messages.")
		}
		slog.Info(fmt.Sprintf("Received a message of size %d bytes.", size))
		slog.Debug(string(buffer[:size]))
	}
}