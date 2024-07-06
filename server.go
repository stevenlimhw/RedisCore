package main

import (
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
}

func NewServer(addr string) *Server {
	return &Server{
		listenAddr: addr,
		peers: make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh: make(chan struct{}),
	}
}


func (server *Server) Start() error {
	ln, err := net.Listen("tcp", server.listenAddr)
	if err != nil {
		slog.Error(fmt.Sprintf("Cannot start TCP server at %s", server.listenAddr))
	}
	server.listener = ln

	slog.Info("Checking current peer connections...")
	go server.checkPeerConnections()

	slog.Info(fmt.Sprintf("Starting TCP server at port %s ...", server.listenAddr))
	return server.acceptPeerConnections()
}

// checkPeerConnections handles cases when adding a new peer and
// quitting the channel
func (server *Server) checkPeerConnections() {
	for {
		select {
		case <-server.quitCh:
			slog.Info("Quitting channel.")
			return
		case peer := <-server.addPeerCh:
			slog.Info("Adding channel.")
			server.peers[peer] = true
		}
	}
}

// acceptPeerConnections creates a new listener for any new peer connections to the server
// and handles the reading of messages for the peer connections.
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

// handlePeerConnection reads messages sent to the peer connection.
func (server *Server) handlePeerConnection(conn net.Conn) {
	peer := &Peer{
		conn: conn,
	}
	server.addPeerCh <-peer
	slog.Info("Successfully accepted peer connection.", "remoteAddr", conn.RemoteAddr())

	peer.readMessages()
}