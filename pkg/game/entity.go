package game

// Entity - Game component that occupies a Tile
type Entity struct {
	Card *Card
	BaseAttack int
	BaseHealth int
	Attack int
	Health int
	IsGeneral bool
	IsExhausted bool
}
