package geom

var (
	UP    = Delta{0, -1}
	LEFT  = Delta{-1, 0}
	DOWN  = Delta{0, 1}
	RIGHT = Delta{1, 0}
)

type Delta struct {
	X, Y int
}

func (d Delta) Scale(s int) Delta {
	return Delta{d.X * s, d.Y * s}
}

func (d Delta) TurnLeft() Delta {
	return Delta{d.Y, -d.X}
}

func (d Delta) TurnRight() Delta {
	return Delta{-d.Y, d.X}
}

func (d Delta) Reverse() Delta {
	return Delta{-d.X, -d.Y}
}

func (d Delta) IsRight() bool {
	return d.X > 0
}

func (d Delta) IsLeft() bool {
	return d.X < 0
}

func (d Delta) IsHorizontal() bool {
	return d.X != 0
}

func (d Delta) IsUp() bool {
	return d.Y < 0
}

func (d Delta) IsDown() bool {
	return d.Y > 0
}

func (d Delta) IsVertical() bool {
	return d.Y != 0
}
