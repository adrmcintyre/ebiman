package color

// ladder simulates the resistor ladder used for converting
// the asserted logic levels to currents to drive the display
// circuitry. The result is the proportion of maximum drive
// current, between 0.0 and 1.0
func ladder(bits byte) float64 {
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

// toRGB converts a colorByte to float64 RGB.
// The result is written to the provided pointers.
func (bbgggrrr colorByte) toRGB() (float64, float64, float64) {
	bb := (bbgggrrr >> 6) & 0b11
	bbb := byte(bb << 1)
	ggg := byte((bbgggrrr >> 3) & 0b111)
	rrr := byte((bbgggrrr >> 0) & 0b111)

	return ladder(rrr), ladder(ggg), ladder(bbb)
}
