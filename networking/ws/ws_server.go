package ws

import (
	"fmt"
	"net/http"

	clipboard "github.com/benn0x/clipboard-sync/clipboard/write"
	"github.com/benn0x/clipboard-sync/security"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
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
		case message := <-h.broadcast:
			security.Decrypt <- message
			decryptedMessage := <-security.Decrypt

			clipboard.ClipboardWrite(string(decryptedMessage))
			for client := range h.clients {
				select {
				case client.send <- []byte(message):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func StartWebsocketServer() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(WS_ADDR, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("ListenAndServer: %v", err))
	}
}
