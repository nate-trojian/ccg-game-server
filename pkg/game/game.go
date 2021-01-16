package game

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nate-trojian/ccg-game-server/internal"
	"github.com/nate-trojian/ccg-game-server/pkg/matchmaking"
	"go.uber.org/zap"
)

// Rules - Game rules
type Rules struct {
	MulliganHandSize int
	MulligansAllowed int
	PlayerHandSize int
	MaximumMana int
	BoardTemplate BoardTemplate
}

// Game - It's the Game
type Game struct {
	ID string
	logger *zap.Logger
	Rules Rules
	Player1 *Player
	Player2 *Player
	startTime int64
	Turn int
	ActionChan chan Action
	Player1OutChan chan []byte
	Player2OutChan chan []byte
	Board *Board
	Hooks []Hook
}

// NewGame creates a new Game
func NewGame(db Database, match matchmaking.Match, p1Chan, p2Chan chan []byte) (*Game, error) {
	id := uuid.New().String()
	logger := internal.NewLogger(id)

	p1Info, err := db.PlayerInfoFromId(match.Player1)
	if err != nil {
		return nil, fmt.Errorf("failed to get player 1 info - %w", err)
	}

	p1DeckInfo, err := DecodeBase64(match.Player1Deck)
	if err != nil {
		return nil, fmt.Errorf("failed to decode player 1 deck - %w", err)
	}

	p2Info, err := db.PlayerInfoFromId(match.Player2)
	if err != nil {
		return nil, fmt.Errorf("failed to get player 2 info - %w", err)
	}

	p2DeckInfo, err := DecodeBase64(match.Player2Deck)
	if err != nil {
		return nil, fmt.Errorf("failed to decode player 2 deck - %w", err)
	}

	return &Game{
		ID: id,
		logger: logger,
		Player1: &Player{
			Info: p1Info,
			Deck: createDeckFromInfo(p1DeckInfo),
		},
		Player2: &Player{
			Info: p2Info,
			Deck: createDeckFromInfo(p2DeckInfo),
		},
		Rules: getRulesFromMode(match.Mode),
		ActionChan: make(chan Action, 10),
		Player1OutChan: p1Chan,
		Player2OutChan: p2Chan,
	}, nil
}

func getRulesFromMode(mode string) Rules {
	return Rules{
		MulliganHandSize: 4,
		MulligansAllowed: 4,
		PlayerHandSize: 10,
		MaximumMana: 10,
		BoardTemplate: BoardTemplate{},
	}
}

func createDeckFromInfo(info *DeckInfo) *Deck {
	return nil
}

// GetPlayer - Get player by index
func (g *Game) GetPlayer(n int) *Player {
	switch n {
	case 1: return g.Player1
	case 2: return g.Player2
	default: return nil
	}
}

// GetPlayerFromID - Get player with id
func (g *Game) GetPlayerFromID(id string) *Player {
	switch id {
	case g.Player1.Info.ID: return g.Player1
	case g.Player2.Info.ID: return g.Player2
	default: return nil
	}
}

// Start Game process
func (g *Game) Start() {
	g.startTime = time.Now().Unix()
	g.Turn = 0
	g.initializeBoard()
	g.Player1.Deck.Shuffle()
	g.Player2.Deck.Shuffle()
}

func (g *Game) initializeBoard() {
	var (
		ind int
	)

	b := &Board{
		width: g.Rules.BoardTemplate.Width,
		height: g.Rules.BoardTemplate.Height,
		tiles: make([]Tile, g.Rules.BoardTemplate.Width * g.Rules.BoardTemplate.Height),
	}

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			ind = j + i * b.width
			b.tiles[ind] = Tile{
				Entity: g.Rules.BoardTemplate.Entities[ind],
				TileEffect: g.Rules.BoardTemplate.TileEffects[ind],
			}
		}
	}

	for p, pos := range(g.Rules.BoardTemplate.Generals) {
		b.tiles[pos].Entity = g.GetPlayer(p).General
	}

	g.Board = b
}

func (g *Game) mulligan() {
	for i := 0; i < g.Rules.MulliganHandSize; i++ {
	}
}

func (g *Game) loop() {
}

func (g *Game) processAction(a Action) {
	
}