package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
)

type Server struct {
	listener   net.Listener
	addPeerCh  chan *Peer
	messageCh  chan *Message
	peers      map[*Peer]bool
	quitCh     chan struct{}
	listenAddr string
}

func NewServer(addr string) *Server {
	return &Server{
		addPeerCh:  make(chan *Peer),
		messageCh:  make(chan *Message),
		peers:      make(map[*Peer]bool),
		listenAddr: addr,
		quitCh:     make(chan struct{}),
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

	value := message.value
	if value.typ != "array" {
		slog.Error("Invalid request: expected RESP array.")
		return
	}
	if len(value.array) == 0 {
		slog.Error("Invalid request: expected non-empty RESP array.")
	}

	command := strings.ToUpper(value.array[0].bulk)
	args := value.array[1:]

	writer := NewWriter(message.peer.conn)

	handler, ok := Handlers[command]
	if !ok {
		slog.Debug("Invalid command inputted:", "command", command)
		err := writer.Write(&Value{typ: "string", str: ""})
		if err != nil {
			slog.Error("Failed to write value to connection.")
		}
		return
	}

	result := handler(args)
	err := writer.Write(result)
	if err != nil {
		slog.Error("Failed to write value to connection.")
	}

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
		conn:      conn,
		quitCh:    server.quitCh,
		messageCh: server.messageCh,
	}
	// add peer connection to server via channel
	server.addPeerCh <- peer
	slog.Info("Successfully accepted peer connection.", "remoteAddr", conn.RemoteAddr())

	peer.readMessages()
}
