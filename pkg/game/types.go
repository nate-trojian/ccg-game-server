package game

// Player - 
type Player struct {
	ID string
	Username string
}

// Game - 
type Game struct {
	ID string
	Player1 Player
	Player2 Player
	Player1Deck Deck
	Player2Deck Deck
}

// EffectFunc - 
type EffectFunc func(Card)

// Effect - 
type Effect struct {
	ID string
	Modifier EffectFunc
}

// Card - 
type Card struct {
	ID string
	ImageID string
	Name string
	ManaCost int
	Text string
	Effects []Effect
}

// Creature - 
type Creature struct {
	Card
	BaseAttack int
	BaseHealth int
	Attack int
	Health int
}

// Spell - 
type Spell struct {
	Card
}

// Deck -
type Deck struct {
	ID string
	Cards []Card
}