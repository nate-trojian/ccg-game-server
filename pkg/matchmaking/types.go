package matchmaking

// Request - Client Request to be put in Matchmaking
type Request struct {
	PlayerID string `json:"P"`
	Mode     string `json:"M"`
	DeckID   string `json:"D"`
}

// Match - Details to start a match
type Match struct {
	Player1     string `json:"P1"`
	Player1Deck string `json:"D1"`
	Player2     string `json:"P2"`
	Player2Deck string `json:"D2"`
	Mode        string `json:"M"`
	ID          string `json:"ID"`
}
