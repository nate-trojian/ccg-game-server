package pkg

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nate-trojian/ccg-game-server/internal"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// Send channel buffer
	sendChannelBuffer = 0  // Will need to test to figure out where this should be at
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client - Websocket Client
// Handles all reading and writing operations
type Client struct {
	logger *zap.Logger
	conn *websocket.Conn
	Send chan []byte
}

// NewClient - Creates a new Client from a websocket connection
func NewClient(reqAddr string, conn *websocket.Conn) *Client {
	return &Client{
		logger: internal.NewLogger("client").With(zap.String("ip", reqAddr)),
		conn: conn,
		Send: make(chan []byte, sendChannelBuffer),
	}
}

func (c *Client) read(msg chan []byte, close chan *Client) {
	defer func() {
		close <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("Websocket unexpectedly closed", zap.Error(err))
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// Pass message to hub so it gets passed to the right game
		msg <- message
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.logger.Error("Failed to create new writer", zap.Error(err))
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				c.logger.Error("Failed to close writer", zap.Error(err))
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}