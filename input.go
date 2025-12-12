package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func GetJoystickInput() int {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		return data.JOY_LEFT
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		return data.JOY_RIGHT
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		return data.JOY_UP
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		return data.JOY_DOWN
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		return data.JOY_BUTTON
	}
	return data.JOY_NONE
}
