package websocket

import (
	"sync"
)

type Message struct {
	RoomID  uint
	Content []byte
	Client  *Client
}

type Hub struct {
	Rooms      map[uint]map[*Client]bool
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[uint]map[*Client]bool),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.handleRegister(client)

		case client := <-hub.Unregister:
			hub.handleUnregister(client)

		case message := <-hub.Broadcast:
			hub.handleBroadcast(message)
		}
	}
}

func (hub *Hub) handleRegister(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	room, ok := hub.Rooms[client.RoomID]
	if !ok {
		room = make(map[*Client]bool)
		hub.Rooms[client.RoomID] = room
	}
	room[client] = true
}

func (hub *Hub) handleUnregister(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if room, ok := hub.Rooms[client.RoomID]; ok {
		if _, exist := room[client]; exist {
			delete(room, client)
			close(client.Send)
		}
		if len(room) == 0 {
			delete(hub.Rooms, client.RoomID)
		}
	}
}

func (hub *Hub) handleBroadcast(message *Message) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if room, ok := hub.Rooms[message.RoomID]; ok {
		var toRemove []*Client
		for client := range room {
			select {
			case client.Send <- message.Content:
			default:
				toRemove = append(toRemove, client)
				close(client.Send)
			}
		}
		for _, client := range toRemove {
			delete(room, client)
		}
		if len(room) == 0 {
			delete(hub.Rooms, message.RoomID)
		}
	}
}
