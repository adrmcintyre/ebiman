package input

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	JOY_NONE   = 0
	JOY_UP     = 1
	JOY_LEFT   = 2
	JOY_DOWN   = 4
	JOY_RIGHT  = 8
	JOY_CENTRE = 16
	JOY_BUTTON = 32
)

var JoyDirection = map[int]geom.Delta{
	JOY_UP:    geom.UP,
	JOY_LEFT:  geom.LEFT,
	JOY_DOWN:  geom.DOWN,
	JOY_RIGHT: geom.RIGHT,
}

func GetJoystickSwitch() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

// TODO - monitor key presses and releases so the most recently pressed key
// always takes precedence when multiple keys are down
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
