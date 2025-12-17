package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/video"
)

const (
	BONUS_X = PACMAN_START_X
	BONUS_Y = 160
)

type BonusActor struct {
	Visible bool
	Pos     geom.Position
}

func MakeBonus() BonusActor {
	return BonusActor{
		Pos: geom.Position{BONUS_X, BONUS_Y},
	}
}

func (b *BonusActor) BonusInit() {
	b.Visible = false
}

func (b *BonusActor) BonusStart() {
	b.Visible = false
}

func (b *BonusActor) DrawBonus(v *video.Video, bonusInfo data.BonusInfoEntry) {
	look := bonusInfo.Sprite
	pal := bonusInfo.Pal
	if b.Visible {
		offset := geom.Delta{-4, -4 - MAZE_TOP}
		v.AddSprite(b.Pos.Add(offset), look, pal)
	}
}
