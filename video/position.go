package video

type ScreenPos struct {
	X int
	Y int
}

func (p ScreenPos) ToTilePos() TilePos {
	return TilePos{p.X / 8, p.Y / 8}
}

type TilePos struct {
	X int
	Y int
}

func (pos TilePos) ToScreenPos() ScreenPos {
	return ScreenPos{pos.X * 8, pos.Y * 8}
}

func (p TilePos) DistSq(q TilePos) int {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return dx*dx + dy*dy
}

// Tile co-ord conversion:
//
//	top (0 <= x < 32, y < 2) - note x=0,1,30,31 are invisible
//	index := (y+30)*32 + (31-x)				// 0x3c0-0x3ff
//
//	normal (0 <= x < 28, 2 <= y < 34)
//	index := (29-pos.X)*32 + (pos.Y - 2) 	// 0x040-0x3bf
//
//	bottom (0 <= x < 32, y < 2) - note x=0,1,30,31 are invisible
//	index := y*32 + (31-x) 					// 0x000-0x03f
func (pos TilePos) tileIndex() int {
	switch {
	case pos.Y < 2:
		return 0x3c0 + pos.Y*32 + 31 - pos.X
	case pos.Y >= 34:
		return 0x000 + (pos.Y-34)*32 + 31 - pos.X
	default:
		return 0x40 + (27-pos.X)*32 + (pos.Y - 2)
	}
}
