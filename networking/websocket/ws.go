package websocket

import "time"

const (
	WEBSOCKET_ADDR       = ":8080"                        // The address to host websocket on
	WEBSOCKET_WRITE_WAIT = 10 * time.Second               // The allowed time to write a message to the peer
	WEBSOCKET_PONG_WAIT  = 60 * time.Second               // The allowed time to read the next pong message from the peer
	WEBSOCKET_PING_WAIT  = (WEBSOCKET_PONG_WAIT * 9) / 10 // The time between pings, must be lower than the pong wait.
	WEBSOCKET_MAX_SIZE   = 512                            // Maximum message size allowed from a peer
)
