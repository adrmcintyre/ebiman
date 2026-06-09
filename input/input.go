package input

import (
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// A bitmap of joystick buttons
type Joystick int

const (
	JoyNone = Joystick(0)
	JoyUp   = Joystick(1 << iota)
	JoyLeft
	JoyDown
	JoyRight
	JoyButton
)

// JoyDirection maps a joystick input to a heading.
var joyDirection = map[Joystick]geom.Delta{
	JoyUp:    geom.Up,
	JoyLeft:  geom.Left,
	JoyDown:  geom.Down,
	JoyRight: geom.Right,
}

// Direction returns a delta corresponding to the stick state.
func (j Joystick) Direction() (geom.Delta, bool) {
	delta, ok := joyDirection[j]
	return delta, ok
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
func (i *Input) JoystickDirection() Joystick {
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
func (i *Input) JoystickInput() Joystick {
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

// Pause returns true if the pause key has just been pressed.
func (i *Input) Pause() bool {
	return i.IsJustPressed(ebiten.KeyS)
}

// VolumeUp returns true if the volume-up key has just been pressed.
func (i *Input) VolumeUp() bool {
	return i.IsJustPressed(ebiten.KeyP)
}

// VolumeDown returns true if the volume-up key has just been pressed.
func (i *Input) VolumeDown() bool {
	return i.IsJustPressed(ebiten.KeyO)
}
