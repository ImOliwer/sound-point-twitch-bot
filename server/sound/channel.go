package sound

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/imoliwer/sound-point-twitch-bot/server/twitch_irc"
)

type SoundDeployment struct {
	State *twitch_irc.UserState `json:"userstate"`
	Price uint64                `json:"price"`
	Id    string                `json:"id"`
}

type Server struct {
	upgrader websocket.Upgrader
	closed   bool
	clients  map[*websocket.Conn]bool
}

func NewServer(readBuffer int, writebuffer int) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  readBuffer,
			WriteBufferSize: writebuffer,
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

func (r *Server) Listen(address string) {
	r.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		socket, err := r.upgrader.Upgrade(w, req, nil)
		if err != nil {
			panic(err)
		}
		r.clients[socket] = true
	})

	go func() {
		if err := http.ListenAndServe(address, nil); err != nil {
			panic(err)
		}
	}()
}

func (r *Server) Stop() {
	r.closed = true
}

func (r *Server) Broadcast(obj SoundDeployment) {
	for client := range r.clients {
		client.WriteJSON(obj)
	}
}

func (r *Server) register(conn *websocket.Conn) {
	go func() {
		for {
			if !r.closed {
				if _, _, err := conn.ReadMessage(); err == nil {
					continue
				}
			}
			// ensure the closing of said client AND its (this) loop
			r.unregister(conn)
			return
		}
	}()
}

func (r *Server) unregister(conn *websocket.Conn) {
	delete(r.clients, conn)
	conn.Close()
}
