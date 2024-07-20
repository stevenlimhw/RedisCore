package main

import (
	"fmt"
	"log/slog"
	"net"
)

type Peer struct {
	conn net.Conn
  quitCh chan struct{}
  messageCh chan *Message
}

type Message struct {
  peer *Peer 
  value Value
}

func (peer *Peer) readMessages() {
  slog.Info("Reading messages from the peer connection...")
  resp := NewResp(peer.conn)

  for {
    value, err := resp.Parse()
    if err != nil {
      slog.Error("Unable to parse RESP command.")
      peer.quitCh <-struct{}{}
    }
    fmt.Println("PEER: ", value)
    message := &Message{
      peer: peer,
      value: value,
    } 
    peer.messageCh <-message
  }
}

