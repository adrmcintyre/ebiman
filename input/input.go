package input

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TODO create a type
// A bitmap of joystick directions
const (
	JOY_NONE   = 0
	JOY_UP     = 1
	JOY_LEFT   = 2
	JOY_DOWN   = 4
	JOY_RIGHT  = 8
	JOY_CENTRE = 16
	JOY_BUTTON = 32
)

// JoyDirection maps a joystick input to a heading.
var JoyDirection = map[int]geom.Delta{
	JOY_UP:    geom.UP,
	JOY_LEFT:  geom.LEFT,
	JOY_DOWN:  geom.DOWN,
	JOY_RIGHT: geom.RIGHT,
}

// GetJoystickSwitch returns true if the "switch" is currently pressed.
//
// We use the spacebar as a proxy for the button.
func GetJoystickSwitch() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

// GetJoystickDirection describes the current direction of the joystick.
// We use the arrow keys as proxy for the joystick.
func GetJoystickDirection() int {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		return JOY_UP
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		return JOY_LEFT
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		return JOY_DOWN
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		return JOY_RIGHT
	}
	return JOY_CENTRE
}

// GetJoystickInput describes the most recent input from the joystick.
func GetJoystickInput() int {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		return JOY_UP
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		return JOY_LEFT
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		return JOY_DOWN
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		return JOY_RIGHT
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		return JOY_BUTTON
	}
	return JOY_NONE
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
