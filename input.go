package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/hajimehoshi/ebiten/v2"
)

func GetJoystickSwitch() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

// TODO - monitor key presses and releases so the most recently pressed key
// always takes precedence when multiple keys are down
func GetJoystickDirection() int {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		return data.JOY_LEFT
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		return data.JOY_RIGHT
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		return data.JOY_UP
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		return data.JOY_DOWN
	}
	return data.JOY_CENTRE
}
