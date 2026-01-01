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

// TODO - presumably these are screen co-ordinates?
// Move play area to top or left of screen on mobile,
// and render touch targets.

// JoystickSwitch returns true if the "switch" is currently pressed.
//
// We use the spacebar as a proxy for the button.
func JoystickSwitch() bool {
	ids := inpututil.JustPressedTouchIDs()
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		_ = x
		if y < 350 {
			return true
		}
	}
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

// JoystickDirection describes the current direction of the joystick.
// We use the arrow keys as proxy for the joystick.
func JoystickDirection() int {
	ids := inpututil.JustPressedTouchIDs()
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		switch {
		case y >= 350 && x < 100:
			return JoyLeft
		case y >= 350 && x >= 200:
			return JoyRight
		case y >= 350 && y < 400 && x >= 100 && x < 200:
			return JoyUp
		case y >= 400 && x >= 100 && x < 200:
			return JoyDown
		}
	}
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
	ids := inpututil.JustPressedTouchIDs()
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		_ = x
		if y < 350 {
			return JoyButton
		}
		switch {
		case y >= 350 && x < 100:
			return JoyLeft
		case y >= 350 && x >= 200:
			return JoyRight
		case y >= 350 && y < 400 && x >= 100 && x < 200:
			return JoyUp
		case y >= 400 && x >= 100 && x < 200:
			return JoyDown
		}
	}
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
