package main

import (
	"log"

	"github.com/adrmcintyre/ebiman/game"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	BUILD_TAG     = "local"
	NAKAMA_URL    = "http://127.0.0.1:7350"
	NAKAMA_KEY    = "defaultkey"
	IS_WASM_BUILD = "0"
)

// setWindowSize sizes the containing window ratio to have the given
// aspectRatio and fill the width or height of the screen to the given
// fillRatio.
func setWindowSize(aspectRatio float64, fillRatio float64) {
	w, h := ebiten.Monitor().Size()
	fw, fh := float64(w), float64(h)
	if fw/fh > aspectRatio {
		w = int(fh * aspectRatio)
	} else {
		h = int(fw / aspectRatio)
	}
	ebiten.SetWindowSize(int(float64(w)*fillRatio), int(float64(h)*fillRatio))
}

func main() {
	ebiten.SetWindowTitle("ebiman")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	setWindowSize(28.0/36.0, 0.75)

	g := game.NewGame(NAKAMA_URL, NAKAMA_KEY, IS_WASM_BUILD != "0")
	if err := g.Execute(); err != nil {
		log.Fatal(err)
	}
}
