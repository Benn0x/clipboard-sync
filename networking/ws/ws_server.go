package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

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

// Creates and returns a new Hub
func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Starts the run loop
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

// Start monitoring the clipboard (host only)
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

type ipEntry struct {
	IpAddress string
	Count     int
}

// Start a websocket server and have it run on the host
func StartWebsocketServer() {
	hub := newHub()
	go hub.run()
	go StartCPMonitorHost(hub)
	ipConnectionLogs := []ipEntry{}
	ticker := time.NewTicker(CONN_RATE_LIMIT_RESET)

	go func() {
		for range ticker.C {
			ipConnectionLogs = []ipEntry{}
		}
	}()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		userIp, _, err := net.SplitHostPort(r.RemoteAddr)
		found := false
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			fmt.Printf("Rate limit error: %v\n", err)
			return
		}

		if len(authHeader) > 0 {
			security.Decrypt <- []byte(authHeader)
			decryptedMessage := <-security.Decrypt
			if string(decryptedMessage) == security.FAILED_DECRYPTION_MESSAGE {
				http.Error(w, "Unauthorized", 401)
				fmt.Printf("Unauthorized connection attempt from %s\n", userIp)
				return
			}

			userCount := 1
			for i, entry := range ipConnectionLogs {
				if entry.IpAddress == userIp {
					ipConnectionLogs[i].Count += 1
					userCount = ipConnectionLogs[i].Count
					found = true
				}
			}

			if !found {
				ipConnectionLogs = append(ipConnectionLogs, ipEntry{
					IpAddress: userIp,
					Count:     1,
				})
			}

			if userCount > CONN_RATE_LIMIT {
				http.Error(w, "Too Many Requests", 429)
				fmt.Printf("Rate limit exceeded for IP: %s\n", userIp)
				return
			}

			serveWs(hub, w, r)
		} else {
			http.Error(w, "Unauthorized", 401)
			fmt.Printf("Unauthorized connection attempt from %s\n", userIp)
			return
		}
	})
	err := http.ListenAndServe(WS_ADDR, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("ListenAndServer: %v", err))
	}
}
