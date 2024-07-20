package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
)

type Server struct {
	listenAddr string
	listener net.Listener
	peers map[*Peer]bool
	addPeerCh chan *Peer
	quitCh chan struct{}
  messageCh chan *Message
}

func NewServer(addr string) *Server {
	return &Server{
		listenAddr: addr,
		peers: make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh: make(chan struct{}),
    messageCh: make(chan *Message),
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

func (server *Server) handleMessage(message *Message) {
  fmt.Println("SERVER: ", message.value)
  message.peer.conn.Write([]byte("+OK\r\n"))
}

// Handles data coming into the channels in the server.
func (server *Server) handleChannels() {
	for {
		select {
		case <-server.quitCh:
			slog.Info("Received signal to close server. Server closing...")
      server.listener.Close()
      os.Exit(1)
		case peer := <-server.addPeerCh:
			slog.Info("Adding and tracking peer connection to server.")
			server.peers[peer] = true
    case message := <-server.messageCh:
      slog.Info("Received message from message channel.")
      server.handleMessage(message)
		}
	}
}

// Creates a new listener for any new peer connections to the server.
func (server *Server) acceptPeerConnections() error {
  retries := 0
	for {
		conn, err := server.listener.Accept()

    if retries >= 3 && err != nil {
      slog.Error("Failed to connect to peer after retries. Closing connection.")
      return err
    }

		if err != nil {
			slog.Error("Failed to connect to peer. Retrying...")
      retries++
			continue
		}
    slog.Info("Connecting to peer...")
		go server.handlePeerConnection(conn)
	}
}

// Adds a new peer connection to the server and reads messages sent by the peer.
func (server *Server) handlePeerConnection(conn net.Conn) {
	peer := &Peer{
		conn: conn,
    quitCh: server.quitCh,
    messageCh: server.messageCh,
	}
  // add peer connection to server via channel
	server.addPeerCh <-peer
	slog.Info("Successfully accepted peer connection.", "remoteAddr", conn.RemoteAddr())
  
	peer.readMessages()
}
