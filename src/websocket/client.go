package websocket

import (
	application_interfaces "first-project/src/application/interfaces"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeTimeout      = 10 * time.Second
	readTimeout       = 60 * time.Second
	pingPeriod        = (readTimeout * 9) / 10
	messageBufferSize = 256
)

type Client struct {
	Hub         *Hub
	Conn        *websocket.Conn
	Send        chan []byte
	RoomID      uint
	UserID      uint
	mu          sync.Mutex
	done        chan struct{}
	chatService application_interfaces.ChatService
}

func NewClient(hub *Hub, conn any, roomID, userID uint, chatService application_interfaces.ChatService) *Client {
	if hub == nil {
		panic("hub cannot be nil")
	}
	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		panic("invalid connection type")
	}
	return &Client{
		Hub:         hub,
		Conn:        wsConn,
		Send:        make(chan []byte, messageBufferSize),
		RoomID:      roomID,
		UserID:      userID,
		done:        make(chan struct{}),
		chatService: chatService,
	}
}

func (client *Client) ReadPump() error {
	defer func() {
		client.Hub.Unregister <- client
		close(client.done)
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(readTimeout))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			return err
		}
		client.chatService.SaveMessage(client.RoomID, client.UserID, string(message))
		client.Hub.Broadcast <- &Message{
			RoomID:  client.RoomID,
			Content: message,
			Client:  client,
		}
	}
}

func (client *Client) WritePump() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.mu.Lock()
			client.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				client.mu.Unlock()
				return fmt.Errorf("send channel closed")
			}
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			client.mu.Unlock()
			if err != nil {
				return err
			}
		case <-ticker.C:
			client.mu.Lock()
			client.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			err := client.Conn.WriteMessage(websocket.PingMessage, nil)
			client.mu.Unlock()
			if err != nil {
				return err
			}
		}
	}
}
