package game

// PlayerInfo - Player Information 
type PlayerInfo struct {
	ID string
	Username string
}

// Player - Game Player
type Player struct {
	Info *PlayerInfo
	Deck *Deck
	General *Entity
	Mana int
	Health int
}