package matchmaking

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/nate-trojian/ccg-game-server/pkg"
	"go.uber.org/zap"
)

const (
	makerCheckPeriod = 2000
	matchBuffer = 10
)

// Matchmaker - Assigns players to new matches
type Matchmaker struct {
	logger *zap.Logger
	queue map[string][]string
	in chan Request
	out chan Match
}

// InitializeMatchmaker - Creates a new Matchmaker
func InitializeMatchmaker() *Matchmaker {
	return &Matchmaker{
		logger: pkg.NewLogger("matchmaker"),
		queue: map[string][]string{
			"standard": {},
		},
		in: make(chan Request),
		out: make(chan Match, matchBuffer),
	}
}

// In - Returns the intake channel for the Matchmaker
func (m *Matchmaker) In() chan Request {
	return m.in
}

// Out - Returns the output channel for the Matchmaker
func (m *Matchmaker) Out() chan Match {
	return m.out
}

// Start - Starts the matchmaker
func (m *Matchmaker) Start(ctx context.Context) {
	t := time.NewTicker(time.Millisecond * makerCheckPeriod)
	defer func() {
		t.Stop()
		close(m.in)
		close(m.out)
	}()
	for {
		select {
		case req := <-m.in:
			q := m.queue[req.Mode]
			m.queue[req.Mode] = append(q, req.Player)
		case <-t.C:
			rand.Seed(time.Now().UnixNano())
			for k := range m.queue {
				q := m.queue[k]
				// Shuffle the players
				rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })

				// Assign new matches
				for len(q) > 1 {
					p1 := q[0]
					p2 := q[1]
					q = q[2:]
					match := Match{
						Player1: p1,
						Player2: p2,
						Mode: k,
						ID: uuid.New().String(),
					}
					m.logger.Debug("Created new match", zap.Any("match", match))
					m.out <- match
				}
			}
		case <-ctx.Done():
			m.logger.Info("Closing Matchmaker")
			return
		}
	}
}