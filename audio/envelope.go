package audio

// A counter increments a value every 1/60 second.
type counter interface {
	Count() byte
}

// applyEnvelope modulates volume according to the specified envelope.
func applyEnvelope(c counter, vol byte, envelope byte) byte {
	switch envelope {
	case envConst:
		return vol

	case envDecay1, envDecay2, envDecay4, envDecay8:
		rate := byte(1) << (envelope - 1)
		if c.Count()&(rate-1) == 0 {
			if vol&0x0f != 0 {
				return vol - 1
			}
		}
		return vol

	// missing implementation for attack
	default:
		return vol
	}
}
