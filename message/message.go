package message

import (
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/video"
)

type MsgId int

const (
	None MsgId = iota
	Player1
	Player2
	NoPlayer
	Ready
	GameOver
	NoStatus
)

type Message struct {
	X, Y int
	Text string
	Pal  byte
}

var msgs = map[MsgId]Message{
	NoPlayer: {9, 14, "          ", palette.BLACK},
	Player1:  {9, 14, "PLAYER ONE", palette.INKY},
	Player2:  {9, 14, "PLAYER TWO", palette.INKY},
	NoStatus: {9, 20, "          ", palette.BLACK},
	Ready:    {9, 20, "  READY!  ", palette.PACMAN},
	GameOver: {9, 20, "GAME  OVER", palette.PAL_29}, // red
}

func (id MsgId) Draw(v *video.Video) {
	if msg, ok := msgs[id]; ok {
		v.SetCursor(msg.X, msg.Y)
		v.WriteString(msg.Text, msg.Pal)
	}
}
