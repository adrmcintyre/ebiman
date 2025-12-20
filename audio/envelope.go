package audio

// A counter increments a value every 1/60 second.
type counter interface {
	GetCount() byte
}

// applyEnvelope modulates volume according to the specified envelope.
func applyEnvelope(c counter, vol byte, envelope byte) byte {
	switch envelope {
	case ENV_CONST:
		return vol

	case ENV_DECAY1, ENV_DECAY2, ENV_DECAY4, ENV_DECAY8:
		rate := byte(1) << (envelope - 1)
		if c.GetCount()&(rate-1) == 0 {
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
