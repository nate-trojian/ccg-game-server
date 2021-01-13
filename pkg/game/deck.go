package game

import (
	"math/rand"
	"time"
)

// Card -
type Card struct {
	ID string
	Faction int
	OwnedBy PlayerRef
	ImageID string
	Name string
	ManaCost int
	Text string
}

type entityI interface {
	IsGeneral() bool
}

// Entity - Game component that occupies a Tile
type Entity struct {
	Card
	entityI
	BaseAttack int
	BaseHealth int
	Attack int
	Health int
	IsGeneral bool
	IsExhausted bool
}

// Spell - Game component that triggers an Effect
type Spell struct {
	Card
}

// Deck -
type Deck struct {
	ID string
	Cards []Card
	General Entity
}

// Shuffle Deck
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]})
}

// Draw Top Card from Deck
func (d *Deck) Draw() Card {
	top := d.Cards[0]
	d.Cards = d.Cards[1:]
	return top
}

// Replace Card with another Card of a different ID
func (d *Deck) Replace(c Card) Card {
	var (
		j int
		cr Card
	)

	// If deck is empty, return c
	// This way we don't stop someone from replacing even if they have run out of cards
	// I think it's more fun this way
	if len(d.Cards) == 0 {
		return c
	}

	// Get a random start position
	r := rand.Intn(len(d.Cards))

	// Cycle through whole deck
	for i := 0; i < len(d.Cards); i++ {
		// Get usable index
		j = (r+i) % len(d.Cards)
		// Save card
		cr = d.Cards[j]
		// If its a valid replace, exit loop
		if cr.ID != c.ID {
			break
		}
	}
	// Swap card in place
	d.Cards[j] = c
	// Return replaced card
	return cr
}

