package main

import (
	"log/slog"
	"net"
)

type Peer struct {
	conn net.Conn
  messageCh chan []byte
}

func (peer *Peer) readMessages() {
	slog.Info("Reading messages sent to peer...")
	for {
		buffer := make([]byte, 1024)
		size, err := peer.conn.Read(buffer)
		if err != nil {
			slog.Error("Error in reading messages.")
		}

    messageBuffer := make([]byte, size)
    copy(messageBuffer, buffer)
    peer.messageCh <-messageBuffer
    //message := string(buffer[:size])
    //slog.Info(message)
	}
}
