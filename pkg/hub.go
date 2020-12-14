package pkg

import (
	"context"

	"github.com/nate-trojian/ccg-game-server/internal"
	"github.com/nate-trojian/ccg-game-server/pkg/matchmaking"
	"go.uber.org/zap"
)

// Hub - Central delegator of new notifications
type Hub struct {
	// Internal
	logger *zap.Logger
	// Matchmaking
	newMatches chan matchmaking.Match
	// Clients
	register   chan *Client
	unregister chan *Client
	connected  map[*Client]bool
	// Game
	messages chan []byte
}

// NewHub - Create a new Hub instance
func NewHub(nm chan matchmaking.Match) *Hub {
	return &Hub{
		logger:     internal.NewLogger("hub"),
		newMatches: nm,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		connected:  make(map[*Client]bool),
		messages:   make(chan []byte),
	}
}

// RegisterClientsChan - Returns channel to send new clients
func (h *Hub) RegisterClientsChan() chan *Client {
	return h.register
}

// Start - Start the Hub process
func (h *Hub) Start(ctx context.Context) {
	for {
		select {
		case c := <-h.register:
			h.connected[c] = true
			go c.read(h.messages, h.unregister)
		case c := <-h.unregister:
			delete(h.connected, c)
			close(c.Send)
		case m := <-h.newMatches:
			h.logger.Info("Creating new match", zap.Any("match", m))
		}
	}
}
