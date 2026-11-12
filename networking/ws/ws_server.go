package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	clip "github.com/benn0x/clipboard-sync/clipboard/write"
	"github.com/benn0x/clipboard-sync/security"
	"golang.design/x/clipboard"
)

type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	var lastMessage string
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Println("New client connected to the websocket.")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			security.Decrypt <- message
			decryptedMessage := <-security.Decrypt
			if lastMessage == string(decryptedMessage) {
				continue
			}
			fmt.Printf("Updating clipboard contents for %d clients.\n", len(h.clients))

			clip.ClipboardWrite(string(decryptedMessage))
			for client := range h.clients {
				select {
				case client.send <- []byte(message):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

			lastMessage = string(decryptedMessage)
		}
	}
}

func StartCPMonitorHost(h *Hub) {
	ctx := context.Background()

	for data := range clipboard.Watch(ctx) {
		clipboardFormatted := clip.ClipboardContent{
			DataType: int(data.Format),
			Bytes:    data.Bytes,
		}

		bytes, err := json.Marshal(clipboardFormatted)
		if err != nil {
			fmt.Println(fmt.Errorf("Couldn't process clipboard data: %v", err))
		} else {
			security.Encrypt <- string(bytes)
			encryptedClipboard := <-security.Encrypt

			h.Broadcast <- []byte(encryptedClipboard)
		}
	}
}

func StartWebsocketServer() {
	hub := newHub()
	go hub.run()
	go StartCPMonitorHost(hub)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(WS_ADDR, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("ListenAndServer: %v", err))
	}
}
