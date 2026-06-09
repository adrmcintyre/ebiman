package main

import (
	"github.com/adrmcintyre/ebiman/video"
)

// An MsgId identifies a specific message.
type MsgId int

const (
	MsgNone        MsgId = iota // display no message (do not place any tiles)
	MsgPlayer1                  // place "PLAYER ONE" message tiles
	MsgPlayer2                  // place "PLAYER TWO" message tiles
	MsgClearPlayer              // blank out tiles for current player message
	MsgReady                    // place "READY" message tiles
	MsgGameOver                 // place "GAME OVER" message tiles
	MsgClearStatus              // blank out tiles for status message
)

// A message describes how to display a message.
type message struct {
	X, Y int           // location of leftmost tile
	Text string        // the message text
	Pal  video.Palette // how to colour the message
}

// messages defines the message attributes for each Id, except for None.
var messages = map[MsgId]message{
	MsgClearPlayer: {9, 14, "          ", video.PalBlack},
	MsgPlayer1:     {9, 14, "PLAYER ONE", video.PalInky},
	MsgPlayer2:     {9, 14, "PLAYER TWO", video.PalInky},
	MsgClearStatus: {9, 20, "          ", video.PalBlack},
	MsgReady:       {9, 20, "  READY!  ", video.PalPacman},
	MsgGameOver:    {9, 20, "GAME  OVER", video.Pal29}, // red
}

// Draw places the tiles for the identified message.
func (id MsgId) Draw(v *video.Video) {
	if msg, ok := messages[id]; ok {
		v.SetCursor(msg.X, msg.Y)
		v.WriteString(msg.Text, msg.Pal)
	}
}
