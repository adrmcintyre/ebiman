package message

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/video"
)

// An Id identifies a specific message.
type Id int

const (
	None        Id = iota // display no message (do not place any tiles)
	Player1               // place "PLAYER ONE" message tiles
	Player2               // place "PLAYER TWO" message tiles
	ClearPlayer           // blank out tiles for current player message
	Ready                 // place "READY" message tiles
	GameOver              // place "GAME OVER" message tiles
	ClearStatus           // blank out tiles for status message
)

// A message describes how to display a message.
type message struct {
	X, Y int           // location of leftmost tile
	Text string        // the message text
	Pal  color.Palette // how to colour the message
}

// msgs defines the message attributes for each Id, except for None.
var msgs = map[Id]message{
	ClearPlayer: {9, 14, "          ", color.PalBlack},
	Player1:     {9, 14, "PLAYER ONE", color.PalInky},
	Player2:     {9, 14, "PLAYER TWO", color.PalInky},
	ClearStatus: {9, 20, "          ", color.PalBlack},
	Ready:       {9, 20, "  READY!  ", color.PalPacman},
	GameOver:    {9, 20, "GAME  OVER", color.Pal29}, // red
}

// Draw places the tiles for the identified message.
func (id Id) Draw(v *video.Video) {
	if msg, ok := msgs[id]; ok {
		v.SetCursor(msg.X, msg.Y)
		v.WriteString(msg.Text, msg.Pal)
	}
}
