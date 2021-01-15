package game

const (
	// Mulligan - Replace cards in starting hand
	Mulligan = "MULLIGAN"
	// EndTurn - Player Ends turn
	EndTurn = "END_TURN"
	// Surrender - Player Surrenders
	Surrender = "SURRENDER"
	// Replace - Player replaces a card in hand
	Replace = "REPLACE"
	// Play - Play a card from hand
	Play = "PLAY"
	// Move - Move Entity on board
	Move = "MOVE"
	// Attack - Entity attacks another Entity on board
	Attack = "ATTACK"
)

// Action - How the player interacts with the game
type Action struct {
	Type string // One of the constant values above, maybe make this a custom type?
	Who string // Player ID - Make this a PlayerRef?
	When int64  // Epoch time of the event
	Details map[string]interface{}  // Is it more or less flexible to have the Details map?
}
