package game

// TileEffect -
type TileEffect struct {
	ID string
}

// Tile - 
type Tile struct {
	Entity *Entity
	TileEffect *TileEffect
}

// Occupied - Returns true if tile has reference to an Entity
func (t Tile) Occupied() bool {
	return t.Entity == nil
}

// Board - Game Board
type Board struct {
	width int
	height int
	tiles []Tile
}

// Occupied - Returns true if tile is Occupied
func (b *Board) Occupied(pos int) bool {
	return b.tiles[pos].Occupied()
}

// BoardTemplate - Template that defines an initial board configuration
type BoardTemplate struct {
	Width int
	Height int
	Entities map[int]*Entity
	TileEffects map[int]*TileEffect
	Generals map[int]int  // Player Num -> Tile Pos
}
