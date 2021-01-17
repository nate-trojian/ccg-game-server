package game

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"time"
)

// Skill - JSON representation of what a card does
type Skill struct {
	EffectName string
	Details map[string]interface{}
}

// CardInfo - Base Information about a card
type CardInfo struct {
	ID string
	Faction int
	CardType int
	ImageID string

	Name string
	Mana int
	Text string
	OnPlay []Skill
	OnReplace []Skill
}

// Card - A playable card in the Game
type Card struct {
	CardInfo
	OwnedBy string
}

// Loadout - Deck Loadout
// I'm not sure if this is where this should go, but it works here for now
type Loadout struct {
	General string
}

// DeckInfo - Information to create in-game Deck
type DeckInfo struct {
	Name string
	CardIds []string
	Loadout Loadout
}

// Base64 - Base64 encoded string of the deck
func (di *DeckInfo) Base64() string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	// Ignore the potential error
    _ = json.NewEncoder(encoder).Encode(di)
    encoder.Close()
    return buf.String()
}

// DecodeBase64 - Decode base64 encoded DeckInfo
func DecodeBase64(b64 string) (*DeckInfo, error) {
	buf := bytes.NewBufferString(b64)
	decoder := base64.NewDecoder(base64.StdEncoding, buf)
	var info DeckInfo
	err := json.NewDecoder(decoder).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// Deck - Deck of Cards for Game
type Deck struct {
	Info *DeckInfo
	Cards []*Card
}

// Shuffle Deck
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]})
}

// Draw Top Card from Deck
func (d *Deck) Draw() *Card {
	top := d.Cards[0]
	d.Cards = d.Cards[1:]
	return top
}

// Replace Card with another Card of a different ID
func (d *Deck) Replace(c *Card) *Card {
	var (
		j int
		cr *Card
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

