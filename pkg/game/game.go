package game

import (
	"time"

	"github.com/google/uuid"
	"github.com/nate-trojian/ccg-game-server/pkg/matchmaking"
)

// Rules - Game rules
type Rules struct {
	MulliganHandSize int
	MulligansAllowed int
	PlayerHandSize int
	MaximumMana int
}

// Game - It's the Game
type Game struct {
	ID string
	Rules Rules
	Player1 *Player
	Player2 *Player
	startTime int64
	Turn int
	ActionChan chan Action
	Player1OutChan chan []byte
	Player2OutChan chan []byte
	Template BoardTemplate
	Board *Board
	Hooks []Hook
}

// NewGame creates a new Game
func NewGame(match matchmaking.Match) *Game {
	return &Game{
		ID: uuid.New().String(),
		Player1: &Player{
			PlayerRef: p1,
			Deck: p1Deck,
		},
		Player2: &Player{
			PlayerRef: p2,
			Deck: p2Deck,
		},
		Template: template,
		ActionChan: make(chan Action, 10),
		Player1OutChan: p1Chan,
		Player2OutChan: p2Chan,
	}
}

func getTemplateFromType() BoardTemplate {

}

// GetPlayer - Get player by number
func (g *Game) GetPlayer(n int) *Player {
	switch n {
	case 1: return g.Player1
	case 2: return g.Player2
	default: return nil
	}
}

// GetPlayerFromRef - Get player by ref
func (g *Game) GetPlayerFromRef(ref PlayerRef) *Player {
	switch ref.ID {
	case g.Player1.ID: return g.Player1
	case g.Player2.ID: return g.Player2
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
		width: g.Template.Width,
		height: g.Template.Height,
		tiles: make([]Tile, g.Template.Width * g.Template.Height),
	}

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			ind = j + i * b.width
			b.tiles[ind] = Tile{
				Entity: g.Template.Entities[ind],
				TileEffect: g.Template.TileEffects[ind],
			}
		}
	}

	for p, pos := range(g.Template.Generals) {
		b.tiles[pos].Entity = &(g.GetPlayer(p).Deck.General)
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