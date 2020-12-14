package matchmaking

// Request - Client Request to be put in Matchmaking
type Request struct {
	Player string `json:"C"`
	Mode   string `json:"M"`
}

// Match - Details to start a match
type Match struct {
	Player1 string `json:"P1"`
	Player2 string `json:"P2"`
	Mode    string `json:"M"`
	ID      string `json:"ID"`
}
