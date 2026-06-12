package actor

import (
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/video"
)

func move(pos geom.Position, dir geom.Delta) geom.Position {
	pos = pos.Add(dir)

	// account for tunnel:
	if pos.X <= geom.TunnelLeft && dir.IsLeft() {
		pos.X += geom.TunnelWidth
	} else if pos.X >= geom.TunnelRight && dir.IsRight() {
		pos.X -= geom.TunnelWidth
	}

	return pos
}

func nextTilePos(pos geom.Position, dir geom.Delta) geom.Position {
	return pos.Add(dir.ScaleUp(8)).WrapTunnel()
}

func getTile(v *video.Video, pos geom.Position) video.Tile {
	return v.GetTile(pos.TileXY())
}
