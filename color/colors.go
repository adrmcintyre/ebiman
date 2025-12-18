package color

// simulate resistor ladder
func ladder(bits colorByte) float64 {
	// currents from 5V across 220, 470 and 1k ohm resistors
	const high = 5.0 / 220
	const mid = 5.0 / 470
	const low = 5.0 / 1000
	const limit = (high + mid + low)
	current := 0.0000
	if bits&(1<<2) != 0 {
		current += high
	}
	if bits&(1<<1) != 0 {
		current += mid
	}
	if bits&(1<<0) != 0 {
		current += low
	}
	return current / limit
}

func (bbgggrrr colorByte) toRGB(r, g, b *float64) {
	bb := (bbgggrrr >> 6) & 0b11
	ggg := (bbgggrrr >> 3) & 0b111
	rrr := (bbgggrrr >> 0) & 0b111

	*r = ladder(rrr)
	*g = ladder(ggg)
	*b = ladder(bb << 1) // convert bb to bb0
}
