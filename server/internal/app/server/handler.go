package server

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/stanekondrej/webload/server/pkg/pb/github.com/stanekondrej/webload/protobuf"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		return true // we don't really need this
	},
}

// This type is responsible for handling incoming requests for websocket
// connections.
type Handler struct {
	Conns map[uint32]*struct {
		conn        *websocket.Conn
		lastMessage *protobuf.Stats
	}
}

// Upgrade HTTP connection to WebSockets
func upgradeConnToWs(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	return upgrader.Upgrade(w, r, nil)
}

func (h *Handler) QueryFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := upgradeConnToWs(w, r)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	t, d, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	if t != websocket.TextMessage {
		writeErrorToConn("Expected text message with session ID", conn)
	}

	trimmedData := strings.TrimSpace(string(d))
	i, err := strconv.ParseUint(trimmedData, 10, 32)
	if err != nil {
		log.Println(err)
		return
	}
	id := uint32(i)

	c, ok := h.Conns[id]
	if !ok {
		writeErrorToConn("No such session", conn)
		return
	}

	for {
		if c.lastMessage == nil {
			continue
		}

		bytes, err := proto.Marshal(c.lastMessage)
		if err != nil {
			log.Println(err)
			return
		}

		err = conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func writeErrorToConn(err string, c *websocket.Conn) error {
	return c.WriteMessage(websocket.TextMessage, []byte(err))
}

func (h *Handler) registerNewProvider(c *websocket.Conn) (id uint32) {
	for {
		id = rand.Uint32()
		_, ok := h.Conns[id]

		if ok {
			continue
		}

		h.Conns[id] =
			&struct {
				conn        *websocket.Conn
				lastMessage *protobuf.Stats
			}{
				c,
				nil,
			}
		return
	}
}

func (h *Handler) deregisterProvider(id uint32) {
	delete(h.Conns, id)
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

		h.Conns[id].lastMessage = &status
	}
}

// Construct a new Handler.
func NewHandler() *Handler {
	return &Handler{
		Conns: map[uint32]*struct {
			conn        *websocket.Conn
			lastMessage *protobuf.Stats
		}{},
	}
}
