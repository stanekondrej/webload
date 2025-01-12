package server

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stanekondrej/webload/server/pkg/pb/github.com/stanekondrej/webload/protobuf"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
}

// This type is responsible for handling incoming requests for websocket
// connections.
type Handler struct {
	Ids map[uint32]*websocket.Conn
}

// Upgrade HTTP connection to WebSockets
func upgradeConnToWs(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}

/*
func (h *Handler) QueryFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := upgradeConnToWs(w, r)
	if err != nil {
		log.Println(err)
		return
	}
}
*/

func writeErrorToConn(err string, c *websocket.Conn) error {
	return c.WriteMessage(websocket.TextMessage, []byte(err))
}

func (h *Handler) registerNewProvider(c *websocket.Conn) (id uint32) {
	for {
		id = rand.Uint32()
		_, ok := h.Ids[id]

		if ok {
			continue
		}

		h.Ids[id] = c
		return
	}
}

func (h *Handler) deregisterProvider(id uint32) {
	delete(h.Ids, id)
}

func (h *Handler) ProvideFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := upgradeConnToWs(w, r)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// generate random id and send it over the connection
	id := h.registerNewProvider(conn)
	defer h.deregisterProvider(id)

	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(id)))

	for {
		t, d, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}
		if t != websocket.BinaryMessage {
			log.Println(writeErrorToConn("Invalid message type", conn))
			continue
		}

		var status protobuf.Stats
		err = proto.Unmarshal(d, &status)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println(status)
	}
}

// Construct a new Handler.
func NewHandler() *Handler {
	return &Handler{
		Ids: make(map[uint32]*websocket.Conn),
	}
}
