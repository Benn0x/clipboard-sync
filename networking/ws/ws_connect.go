package ws

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	clipboard "github.com/benn0x/clipboard-sync/clipboard/write"
	"github.com/benn0x/clipboard-sync/security"
	"github.com/gorilla/websocket"
)

var conn *websocket.Conn
var SendOverWS = make(chan string)

func restartAppOnHostClosure(err error) {
	fmt.Printf("Started restart process due to: %v\n", err)

	exe, err := os.Executable()
	if err != nil {
		fmt.Printf("Error during application locating. Error: %v\n", err)
	}

	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error during application restart. Error: %v\n", err)
	}

	os.Exit(0)
}

// Listens for the WS communication.
func ClipboardOverWS() {
	for message := range SendOverWS {
		fmt.Println("Sending clipboard contents over the network.")
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println(fmt.Errorf("WS communicate error: %v", err))
		}
	}
}

// Will connect to the websocket using information given by the MDNS registration
func ConnectToWS(ip string) error {
	url := fmt.Sprintf("ws://%s%s/ws", ip, WS_ADDR)
	dialer := websocket.Dialer{}
	header := make(http.Header)

	security.Encrypt <- rand.Text()
	authorizationHeader := <-security.Encrypt

	header.Set("Authorization", authorizationHeader)

	connection, response, err := dialer.Dial(url, header)
	conn = connection
	if err != nil {
		switch response.StatusCode {
		case http.StatusUnauthorized:
			return fmt.Errorf("Unauthorized access. Please check your passphrase.")
		case http.StatusTooManyRequests:
			return fmt.Errorf("Rate limit exceeded. Please try again in a few minutes.")
		}
		return fmt.Errorf("Dial error: %v", err)
	} else {
		fmt.Printf("Connected to host %s\n", ip)
	}

	go func() {
		go func() {
			ticker := time.NewTicker(PING_PERIOD)
			defer ticker.Stop()

			for range ticker.C {
				err := conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					fmt.Println("Connection lost:", err)
					break
				}
			}
		}()

		for {
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				restartAppOnHostClosure(fmt.Errorf("Host is no longer available to read from."))
				break
			}

			messageType, message, err := conn.ReadMessage()
			if err != nil {
				if closeErr, ok := err.(*websocket.CloseError); ok {
					if closeErr.Code == 1006 {
						restartAppOnHostClosure(fmt.Errorf("Host is no longer available to read from."))
					}
					restartAppOnHostClosure(fmt.Errorf("Connection closed with code %d: %s", closeErr.Code, closeErr.Text))
				}
			} else if messageType == websocket.TextMessage {
				security.Decrypt <- message
				message = <-security.Decrypt
				clipboard.ClipboardWrite(string(message))
			}
		}
	}()

	return nil
}
