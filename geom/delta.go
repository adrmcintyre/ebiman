package geom

var (
	Up    = Delta{0, -1} // heading up the screen
	Left  = Delta{-1, 0} // heading to the left
	Down  = Delta{0, 1}  // heading down the screen
	Right = Delta{1, 0}  // heading to the right
)

// A Delta represents the difference between two Positions.
type Delta struct {
	X, Y int
}

// ScaleUp returns a new Delta scaled up by a factor of s.
func (d Delta) ScaleUp(s int) Delta {
	return Delta{d.X * s, d.Y * s}
}

// TurnLeft rotates the direction vector 90 degrees anti-clockwise.
func (d Delta) TurnLeft() Delta {
	return Delta{d.Y, -d.X}
}

// TurnLeft rotates the direction vector 90 degrees clockwise.
func (d Delta) TurnRight() Delta {
	return Delta{-d.Y, d.X}
}

// TurnLeft rotates the direction vector 180 degrees.
func (d Delta) Reverse() Delta {
	return Delta{-d.X, -d.Y}
}

// IsRight returns true if d represents a rightwards heading.
func (d Delta) IsRight() bool {
	return d.X > 0
}

// IsRight returns true if d represents a leftwards heading.
func (d Delta) IsLeft() bool {
	return d.X < 0
}

// IsRight returns true if d represents a left/right heading.
func (d Delta) IsHorizontal() bool {
	return d.X != 0
}

// IsUp returns true if d represents an upwards heading.
func (d Delta) IsUp() bool {
	return d.Y < 0
}

// IsUp returns true if d represents an downwards heading.
func (d Delta) IsDown() bool {
	return d.Y > 0
}

// IsRight returns true if d represents an up/down heading.
func (d Delta) IsVertical() bool {
	return d.Y != 0
}
