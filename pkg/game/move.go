package game

// Move - How the player interacts with the game
type Move struct {
	Type string
	Who string
	When int64  // Epoch time of the event
}
