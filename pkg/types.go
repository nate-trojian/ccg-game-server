package pkg

// MatchmakingRequest - Client Request to be put in matchmaking
type MatchmakingRequest struct {
	Client string `json:"C"`
	Mode string `json:"M"`
}

// ErrorResponse - Error response message
type ErrorResponse struct {
	Error string `json:"error"`
	Message string `json:"message,omitempty"`
}