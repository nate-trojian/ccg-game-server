package game

// Database - Interface to interact with the database
type Database interface {
	PlayerInfoFromId(string) (*PlayerInfo, error)
	CardInfoFromId(string) (*CardInfo, error)
}