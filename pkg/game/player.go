package game

// PlayerRef -
type PlayerRef struct {
	ID string
	Username string
}

// Player - 
type Player struct {
	PlayerRef
	Deck Deck
	Mana int
	Health int
}