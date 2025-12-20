package message

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/video"
)

type Id int

const (
	None Id = iota
	Player1
	Player2
	ClearPlayer
	Ready
	GameOver
	ClearStatus
)

type message struct {
	X, Y int
	Text string
	Pal  color.Palette
}

var msgs = map[Id]message{
	ClearPlayer: {9, 14, "          ", color.PAL_BLACK},
	Player1:     {9, 14, "PLAYER ONE", color.PAL_INKY},
	Player2:     {9, 14, "PLAYER TWO", color.PAL_INKY},
	ClearStatus: {9, 20, "          ", color.PAL_BLACK},
	Ready:       {9, 20, "  READY!  ", color.PAL_PACMAN},
	GameOver:    {9, 20, "GAME  OVER", color.PAL_29}, // red
}

func (id Id) Draw(v *video.Video) {
	if msg, ok := msgs[id]; ok {
		v.SetCursor(msg.X, msg.Y)
		v.WriteString(msg.Text, msg.Pal)
	}
}
