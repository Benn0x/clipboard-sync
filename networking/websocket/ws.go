package websocket

import "time"

const (
	WEBSOCKET_ADDR       = ":8080"
	WEBSOCKET_WRITE_WAIT = 10 * time.Second
	WEBSOCKET_PONG_WAIT  = 60 * time.Second
	WEBSOCKET_PING_WAIT  = (WEBSOCKET_PONG_WAIT * 9) / 10
	WEBSOCKET_MAX_SIZE   = 512
)
