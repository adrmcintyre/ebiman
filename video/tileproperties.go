package video

// IsTraversable returns true if the tile can be passed over (i.e. not a maze barrier).
func (t Tile) IsTraversable() bool {
	return t == TileSpace || t.IsPill() || t.IsPower() || t >= TileScoreMin && t <= TileScoreMax
}

// IsGate returns true if the tile forms part of the gate guarding the ghosts' home.
func (t Tile) IsGate() bool {
	return t == TileGateLeft || t == TileGateRight
}

// IsHome returns true if the tile is part of the ghosts' home.
func (t Tile) IsHome() bool {
	return t == TileHomeLeft || t == TileHomeRight
}

// IsPill returns true if the tile is a pill of some kind.
func (t Tile) IsPill() bool {
	switch t {
	case TilePill, TilePillMinus, TilePillPlus, TilePillMinus2, TilePillPlus2:
		return true
	default:
		return false
	}
}

// IsPower returns true if the tile is a power pill.
func (t Tile) IsPower() bool {
	switch t {
	case TilePower, TilePowerSmall:
		return true
	default:
		return false
	}
}

// Charge returns the net charge on a tile.
func (t Tile) Charge() int {
	switch t {
	case TilePill:
		return 0
	case TilePillMinus:
		return -1
	case TilePillPlus:
		return 1
	case TilePillMinus2:
		return -2
	case TilePillPlus2:
		return 2
	default:
		return 0
	}
}

// TileFromCharge returns a tile with the specified charge.
func FromCharge(c int) Tile {
	switch c {
	case -1:
		return TilePillMinus
	case 1:
		return TilePillPlus
	case -2:
		return TilePillMinus2
	case 2:
		return TilePillPlus2
	default:
		return TilePill
	}
}
