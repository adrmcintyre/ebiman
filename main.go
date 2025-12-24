package main

import (
	"log"

	"github.com/adrmcintyre/poweraid/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// define a small border between the simulated display and the window border
	hBorder = 8
	vBorder = 8

	// dimensions of simulated display, consisting of 8x8 tiles laid out 28x36.
	hWidth  = 28 * 8
	vHeight = 36 * 8

	// calculate desired physical size of the window
	screenWidth  = hWidth + 2*hBorder
	screenHeight = vHeight + 2*vBorder
	screenScale  = 2.3
)

func main() {
	windowWidth := screenWidth * screenScale
	windowHeight := screenHeight * screenScale

	ebiten.SetWindowTitle("PowerAid")
	ebiten.SetWindowSize(int(windowWidth), int(windowHeight))
	//	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame(screenWidth, screenHeight)
	if err := g.Execute(); err != nil {
		log.Fatal(err)
	}
}
