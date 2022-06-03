package sound

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/imoliwer/sound-point-twitch-bot/server/twitch_irc"
)

type GlobalDeployment struct {
	Price    uint64 `json:"price"`
	ID       string `json:"id"`
	FileName string `json:"file_name"`
}

type RealDeployment struct {
	GlobalDeployment
	State *twitch_irc.UserState `json:"userstate"`
}

type TestDeployment struct {
	GlobalDeployment
	Tester string `json:"tester"`
}

type DeploymentCover struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

func NewCover(readBuffer int, writebuffer int) *DeploymentCover {
	return &DeploymentCover{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  readBuffer,
			WriteBufferSize: writebuffer,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

func (r *DeploymentCover) Broadcast(obj interface{}) {
	for client := range r.clients {
		client.WriteJSON(obj)
	}
}

func (r *DeploymentCover) register(conn *websocket.Conn) {
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err == nil {
				continue
			}
			// ensure the closing of said client AND its (this) loop
			r.unregister(conn)
			return
		}
	}()
}

func (r *DeploymentCover) unregister(conn *websocket.Conn) {
	delete(r.clients, conn)
	conn.Close()
}
