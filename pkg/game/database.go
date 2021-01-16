package game

// Database - Interface to interact with the database
type Database interface {
	PlayerInfoFromId(string) (PlayerInfo, error)
	PlayerDeckFromIds(string, string) (DeckInfo, error)
}