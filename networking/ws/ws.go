package ws

import "time"

const (
	WRITE_WAIT       time.Duration = 10 * time.Second
	PONG_WAIT        time.Duration = 60 * time.Second
	PING_PERIOD      time.Duration = (PONG_WAIT * 9) / 10
	MAX_MESSAGE_SIZE int           = 512

	WS_ADDR string = ":8080"
)

var (
	NEWLINE = []byte{'\n'}
	SPACE   = []byte{' '}
)
