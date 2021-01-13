package game

// Action - Something that directly happens to the game state
type Action struct {
	// Base parameters
	ID string
	game *Game
	parent *Action
	subQueue []Action
	Details map[string]interface{} `json:",omitempty"`

	resolve func()
	sanitize func() Action

	// Targetted Actions
	Source *Entity `json:",omitempty"`
	SourcePos *int `json:",omitempty"`
	Target *Entity `json:",omitempty"`
	TargetPos *int `json:",omitempty"`
}

// SantizedCopy returns a copy of the Action without information for the opposing player or the action itself if it doesn't need to be sanitized
func (a Action) SantizedCopy() Action {
	if a.sanitize != nil {
		return a.sanitize()
	}
	return a
}

// ActionDatabase - Storage for actions
type ActionDatabase map[string]Action
