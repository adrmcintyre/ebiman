package data

const (
	JOY_NONE   = 0
	JOY_UP     = 1
	JOY_LEFT   = 2
	JOY_DOWN   = 4
	JOY_RIGHT  = 8
	JOY_CENTRE = 16
	JOY_BUTTON = 32
)

type Direction struct {
	Dx, Dy int
}

var JoyDirection = map[int]Direction{
	JOY_UP:    {0, -1},
	JOY_LEFT:  {-1, 0},
	JOY_DOWN:  {0, 1},
	JOY_RIGHT: {1, 0},
}
