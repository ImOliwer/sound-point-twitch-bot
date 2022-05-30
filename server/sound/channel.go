package sound

import (
	"bufio"
	"errors"
	"net"
)

type Server struct {
	listening bool
	closed    bool
	clients   map[net.Conn]bool
}

func NewServer() *Server {
	return &Server{
		clients: make(map[net.Conn]bool),
	}
}

func (r *Server) Listen(address string) error {
	if r.listening {
		return errors.New("already listening")
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go func(listener net.Listener, server *Server) {
		for {
			if server.closed {
				listener.Close()
				return
			}

			client, err := listener.Accept()
			if err != nil {
				continue
			}

			server.clients[client] = true
			server.register(client)
		}
	}(listener, r)

	r.listening = true
	return nil
}

func (r *Server) Stop() {
	r.closed = true
}

func (r *Server) Broadcast(message []byte) {
	for client := range r.clients {
		client.Write(message)
	}
}

func (r *Server) register(conn net.Conn) {
	go func() {
		for {
			if _, err := bufio.NewReader(conn).ReadString('\n'); err == nil {
				continue
			}
			// ensure the closing of said client AND its (this) loop
			r.unregister(conn)
			return
		}
	}()
}

func (r *Server) unregister(conn net.Conn) {
	delete(r.clients, conn)
	conn.Close()
}
