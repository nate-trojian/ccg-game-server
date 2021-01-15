package game

import "reflect"

// Effect - Something that directly happens to the game state
type Effect struct {
	// Base parameters
	Type string
	Details map[string]interface{} `json:",omitempty"`
	sensitizeDetails []string

	// Tree parameters
	parent *Effect
	children []Effect
	isDepthFirst bool

	// Execution
	resolve func(*Game)
}

// SantizedCopy returns a copy of the Effect without information for the opposing player or the Effect itself if there's no details to remove
func (e Effect) SantizedCopy() Effect {
	// If there's nothing to do, don't do anything
	if len(e.sensitizeDetails) == 0 {
		return e
	}

	// Create the new copy
	ret := Effect{
		Type: e.Type,
	}

	// Deep copy array
	for k, v := range(e.Details) {
		ret.Details[k] = v
	}

	// Remove everything sensitize
	for _, s := range(e.sensitizeDetails) {
		if _, ok := ret.Details[s]; ok {
			delete(ret.Details, s)
		}
	}

	return ret
}

// Equal - Compares whether an Effect is equal to this one
func (e Effect) Equal(other Effect) bool {
	if e.Type != other.Type {
		return false
	}
	return reflect.DeepEqual(e.Details, other.Details)
}

// InPath returns true if the other Effect is somewhere between this Effect and the root 
func (e Effect) InPath(other Effect) bool {
	// If other is equal to e, we don't have to look any further
	if e.Equal(other) {
		return true
	}
	// Initial set to one up the path
	t := e.parent
	// If t is still valid aka we have a parent to check
	for t != nil {
		// Check if they are equal
		if t.Equal(other) {
			return true
		}
		// If not, go up one in the path
		t = t.parent
	}
	// We hit the root node with no match
	return false
}
