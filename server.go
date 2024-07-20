package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"net"
)

type Server struct {
	listenAddr string
	listener net.Listener
	peers map[*Peer]bool
	addPeerCh chan *Peer
	quitCh chan struct{}
  messageCh chan []byte
}

func NewServer(addr string) *Server {
	return &Server{
		listenAddr: addr,
		peers: make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh: make(chan struct{}),
    messageCh: make(chan []byte),
	}
}


func (server *Server) Start() error {
	ln, err := net.Listen("tcp", server.listenAddr)
	if err != nil {
		slog.Error(fmt.Sprintf("Cannot start TCP server at %s", server.listenAddr))
	}
	server.listener = ln

	slog.Info("Channels are being handled by the server.")
	go server.handleChannels()

	slog.Info(fmt.Sprintf("Starting TCP server at port %s ...", server.listenAddr))
	return server.acceptPeerConnections()
}

func (server *Server) handleMessageBytes(messageBytes []byte) {
  reader := bufio.NewReader(bytes.NewReader(messageBytes))
  resp := &Resp{
    reader: reader,
  }
  resp.parseRespCommand()
}

// Handles data coming into the channels in the server.
func (server *Server) handleChannels() {
	for {
		select {
		case <-server.quitCh:
			slog.Info("Quitting server.")
			return
		case peer := <-server.addPeerCh:
			slog.Info("Adding and tracking peer connection to server.")
			server.peers[peer] = true
    case messageBytes := <-server.messageCh:
      slog.Info("Received message in bytes format.")
      server.handleMessageBytes(messageBytes)
      //slog.Info(string(messageBytes))
		}

	}
}

// Creates a new listener for any new peer connections to the server.
func (server *Server) acceptPeerConnections() error {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			slog.Error("Failed to connect to peer. Retrying...")
			continue
		}
		go server.handlePeerConnection(conn)
	}
}

// Adds a new peer connection to the server and reads messages sent by the peer.
func (server *Server) handlePeerConnection(conn net.Conn) {
	peer := &Peer{
		conn: conn,
    messageCh: server.messageCh,
	}
  // add peer connection to server via channel
	server.addPeerCh <-peer
	slog.Info("Successfully accepted peer connection.", "remoteAddr", conn.RemoteAddr())
  
	peer.readMessages()
}
