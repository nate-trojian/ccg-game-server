package game

// Hook - Hook to resolve after each effect in the Stack
type Hook struct {
	Filter func(Effect) bool
	Effect Effect
}
