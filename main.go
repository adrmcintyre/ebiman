package main

import (
	"log"

	"github.com/adrmcintyre/poweraid/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	hBorder      = 8
	vBorder      = 8
	hWidth       = 28 * 8
	vHeight      = 36 * 8
	screenWidth  = hWidth + 2*hBorder
	screenHeight = vHeight + 2*vBorder
	screenScale  = 2.5
)

func main() {
	ebiten.SetWindowTitle("PowerAid")
	ebiten.SetWindowSize(
		int(screenWidth*screenScale),
		int(screenHeight*screenScale),
	)

	if err := game.EntryPoint(screenWidth, screenHeight); err != nil {
		log.Fatal(err)
	}
}
