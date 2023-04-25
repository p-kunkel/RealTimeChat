package models

import (
	"RealTimeChat/config"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

type Connection struct {
	WS   *websocket.Conn
	Send chan []byte
}

type hub struct {
	Rooms      map[uint64]map[*Connection]bool
	Broadcast  chan Message
	Register   chan Subscription
	Unregister chan Subscription
}

type Subscription struct {
	Conn    *Connection
	Session Session
	ChatId  uint64
}

var ChatHub = hub{
	Broadcast:  make(chan Message),
	Register:   make(chan Subscription),
	Unregister: make(chan Subscription),
	Rooms:      make(map[uint64]map[*Connection]bool),
}

func (s Subscription) ReadPump() {
	c := s.Conn
	m := Message{}
	defer func() {
		ChatHub.Unregister <- s
		c.WS.Close()
	}()

	for {
		_, msg, err := c.WS.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		if err = json.Unmarshal(msg, &m); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		m.ChatId = s.ChatId
		m.SenderId = s.Session.UserId
		m.ReadedBy = pq.Int64Array{int64(s.Session.UserId)}

		if err = m.Create(config.DB); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (s *Subscription) WritePump() {
	c := s.Conn
	ticker := time.NewTicker(time.Second * 5)
	defer func() {
		ticker.Stop()
		c.WS.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Connection) write(mt int, payload []byte) error {
	c.WS.SetWriteDeadline(time.Now().Add(time.Second))
	return c.WS.WriteMessage(mt, payload)
}

func (h *hub) Run() {
	for {
		select {
		case s := <-h.Register:
			connections := h.Rooms[s.ChatId]
			if connections == nil {
				connections = make(map[*Connection]bool)
				h.Rooms[s.ChatId] = connections
			}
			h.Rooms[s.ChatId][s.Conn] = true

		case s := <-h.Unregister:
			connections := h.Rooms[s.ChatId]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)
					close(s.Conn.Send)
					if len(connections) == 0 {
						delete(h.Rooms, s.ChatId)
					}
				}
			}

		case m := <-h.Broadcast:
			connections := h.Rooms[m.ChatId]
			b, _ := json.Marshal(m)

			for c := range connections {
				select {
				case c.Send <- b:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Rooms, m.ChatId)
					}
				}
			}
		}
	}
}
