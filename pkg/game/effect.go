package game

// Effect - Something that directly happens to the game state
type Effect struct {
	// Base parameters
	ID string
	game *Game
	parent *Effect
	subQueue []Effect
	Details map[string]interface{} `json:",omitempty"`

	resolve func()
	sanitize func() Effect

	// Targetted Effects
	Source *Entity `json:",omitempty"`
	SourcePos *int `json:",omitempty"`
	Target *Entity `json:",omitempty"`
	TargetPos *int `json:",omitempty"`
}

// SantizedCopy returns a copy of the Effect without information for the opposing player or the action itself if it doesn't need to be sanitized
func (e Effect) SantizedCopy() Effect {
	if e.sanitize != nil {
		return e.sanitize()
	}
	return e
}
