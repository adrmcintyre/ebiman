package input

import (
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TODO create a type
// A bitmap of joystick directions
const (
	JoyNone   = 0
	JoyUp     = 1
	JoyLeft   = 2
	JoyDown   = 4
	JoyRight  = 8
	JoyCentre = 16
	JoyButton = 32
)

// JoyDirection maps a joystick input to a heading.
var JoyDirection = map[int]geom.Delta{
	JoyUp:    geom.Up,
	JoyLeft:  geom.Left,
	JoyDown:  geom.Down,
	JoyRight: geom.Right,
}

// JoystickSwitch returns true if the "switch" is currently pressed.
//
// We use the spacebar as a proxy for the button.
func JoystickSwitch() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

// JoystickDirection describes the current direction of the joystick.
// We use the arrow keys as proxy for the joystick.
func JoystickDirection() int {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		return JoyUp
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		return JoyLeft
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		return JoyDown
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		return JoyRight
	}
	return JoyCentre
}

// JoystickInput describes the most recent input from the joystick.
func JoystickInput() int {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		return JoyUp
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		return JoyLeft
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		return JoyDown
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		return JoyRight
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		return JoyButton
	}
	return JoyNone
}

// Quit returns true if the quit key has just been pressed.
func Quit() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyQ)
}

func VolumeUp() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyP)
}

func VolumeDown() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyO)
}
