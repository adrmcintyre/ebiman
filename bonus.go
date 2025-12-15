package main

import "github.com/adrmcintyre/poweraid/video"

const (
	BONUS_X = PACMAN_START_X
	BONUS_Y = 160
)

type BonusActor struct {
	Motion Motion
}

func MakeBonus() BonusActor {
	return BonusActor{
		Motion{
			Pos: video.ScreenPos{BONUS_X, BONUS_Y},
		},
	}
}

func (b *BonusActor) BonusInit() {
	b.Motion.Visible = false
}

func (b *BonusActor) BonusStart() {
	b.Motion.Visible = false
}
