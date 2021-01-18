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
	id string
	logger *zap.Logger
	rules Rules
	db Database
	player1 *Player
	player2 *Player
	startTime int64
	turn int
	actionChan chan Action
	player1OutChan chan []byte
	player2OutChan chan []byte
	board *Board
	hooks []Hook
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

	p1Deck, err := createDeckFromInfo(db, match.Player1, p1DeckInfo)
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

	p2Deck, err := createDeckFromInfo(db, match.Player2, p2DeckInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode player 2 deck - %w", err)
	}

	return &Game{
		id: id,
		logger: logger,
		db: db,
		player1: &Player{
			Info: p1Info,
			Deck: p1Deck,
		},
		player2: &Player{
			Info: p2Info,
			Deck: p2Deck,
		},
		rules: getRulesFromMode(match.Mode),
		actionChan: make(chan Action, 10),
		player1OutChan: p1Chan,
		player2OutChan: p2Chan,
	}, nil
}

func getRulesFromMode(mode string) Rules {
	return Rules{
		MulliganHandSize: 4,
		MulligansAllowed: 4,
		PlayerHandSize: 10,
		MaximumMana: 10,
		BoardTemplate: BoardTemplate{
			Width: 9,
			Height: 5,
			Entities: nil,
			TileEffects: map[int]*TileEffect{
				4: {ID: "mana"},
				23: {ID: "mana"},
				41: {ID: "mana"},
			},
			Generals: map[int]int{
				1: 19,
				2: 25,
			},
		},
	}
}

func createDeckFromInfo(db Database, pID string, info *DeckInfo) (*Deck, error) {
	deck := &Deck{
		info,
		[]*Card{},
	}
	for _,c := range(info.CardIds) {
		cInfo, err := db.CardInfoFromId(c)
		if err != nil {
			return nil, err
		}
		deck.Cards = append(deck.Cards, &Card{
			Info: cInfo,
			OwnedBy: pID,
		})
	}
	return deck, nil
}

// GetPlayer - Get player by index
func (g *Game) GetPlayer(n int) *Player {
	switch n {
	case 1: return g.player1
	case 2: return g.player2
	default: return nil
	}
}

// GetPlayerFromID - Get player with id
func (g *Game) GetPlayerFromID(id string) *Player {
	switch id {
	case g.player1.Info.ID: return g.player1
	case g.player2.Info.ID: return g.player2
	default: return nil
	}
}

// Start Game process
func (g *Game) Start() {
	g.startTime = time.Now().Unix()
	g.turn = 0
	g.initializeBoard()
	g.player1.Deck.Shuffle()
	g.player2.Deck.Shuffle()
}

func (g *Game) initializeBoard() {
	var (
		ind int
	)

	b := &Board{
		width: g.rules.BoardTemplate.Width,
		height: g.rules.BoardTemplate.Height,
		tiles: make([]Tile, g.rules.BoardTemplate.Size()),
	}

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			ind = j + i * b.width
			b.tiles[ind] = Tile{
				Entity: g.rules.BoardTemplate.Entities[ind],
				TileEffect: g.rules.BoardTemplate.TileEffects[ind],
			}
		}
	}

	for p, pos := range(g.rules.BoardTemplate.Generals) {
		b.tiles[pos].Entity = g.GetPlayer(p).General
	}

	g.board = b
}

func (g *Game) mulligan() {
	for i := 0; i < g.rules.MulliganHandSize; i++ {
	}
}

func (g *Game) loop() {
}

func (g *Game) processAction(a Action) {
}