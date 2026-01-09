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

// An Input supports user input interactions.
type Input struct {
	touchLayout *TouchLayout
	lastKey     ebiten.Key
}

// New returns an Input.
func New() *Input {
	return &Input{
		touchLayout: NewTouchLayout(),
	}
}

func (i *Input) SetTouchLayout(layout *TouchLayout) {
	i.touchLayout = layout
}

func (i *Input) Update() {
	for _, key := range []ebiten.Key{ebiten.KeySpace, ebiten.KeyUp, ebiten.KeyDown, ebiten.KeyLeft, ebiten.KeyRight} {
		if i.IsJustPressed(key) {
			// TODO - timestamp this and timeout after say 0.25 second?
			// TODO - provide a reset method?
			i.lastKey = key
		}
	}
}

func (i *Input) IsJustPressed(key ebiten.Key) bool {
	return i.IsJustTouched(key) || inpututil.IsKeyJustPressed(key)
}

func (i *Input) IsPressed(key ebiten.Key) bool {
	return i.IsTouched(key) || ebiten.IsKeyPressed(key)
}

// TODO - Move play area to top or left of screen on mobile.
// JoystickSwitch returns true if the "switch" is currently pressed.
//
// We use the spacebar as a proxy for the button.
func (i *Input) JoystickSwitch() bool {
	return i.IsPressed(ebiten.KeySpace)
}

// JoystickDirection describes the current direction of the joystick.
// We use the arrow keys as proxy for the joystick.
func (i *Input) JoystickDirection() int {
	switch i.lastKey {
	case ebiten.KeyUp:
		return JoyUp
	case ebiten.KeyLeft:
		return JoyLeft
	case ebiten.KeyDown:
		return JoyDown
	case ebiten.KeyRight:
		return JoyRight
	}
	return JoyNone
}

// JoystickInput describes the most recent input from the joystick.
func (i *Input) JoystickInput() int {
	switch {
	case i.IsJustPressed(ebiten.KeyUp):
		return JoyUp
	case i.IsJustPressed(ebiten.KeyLeft):
		return JoyLeft
	case i.IsJustPressed(ebiten.KeyDown):
		return JoyDown
	case i.IsJustPressed(ebiten.KeyRight):
		return JoyRight
	case i.IsJustPressed(ebiten.KeySpace):
		return JoyButton
	}
	return JoyNone
}

// Quit returns true if the quit key has just been pressed.
func (i *Input) Quit() bool {
	return i.IsJustPressed(ebiten.KeyQ)
}

func (i *Input) VolumeUp() bool {
	return i.IsJustPressed(ebiten.KeyP)
}

func (i *Input) VolumeDown() bool {
	return i.IsJustPressed(ebiten.KeyO)
}
