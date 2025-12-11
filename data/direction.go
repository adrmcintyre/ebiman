package data

const (
	JOY_UP     = 0
	JOY_LEFT   = 1
	JOY_DOWN   = 2
	JOY_RIGHT  = 3
	JOY_DEAD   = 4
	JOY_CENTRE = 8
	JOY_BUTTON = 16
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
