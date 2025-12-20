package audio

// An EffectProcessor represents the processor for
// a single effects channel.
type EffectProcessor struct {
	counter
	index           int
	command         *Command
	queueMask       uint8
	playingBit      uint8
	envelope        byte
	durationCounter byte
	baseFreq        byte
	vol             byte
	freqDir         bool
	octave          byte
	wave            byte
	initialBaseFreq byte
	freqIncr        byte
	duration        byte
	reverse         bool
	repeatFreqIncr  byte
	repeatCounter   byte
	initialVol      byte
	volIncr         byte
}

var effectChannel [channelCount]EffectProcessor

// processEffects runs all the effects processors
func (au *Audio) processEffects() {
	for _, e := range au.effectProcessor {
		e.command.vol = e.processEffect()
	}
	au.command[0].freq &= 0xffff // retain bottom 16 bits only
}

// clearEffectChannel stops the effect currently in progress
func (e *EffectProcessor) clearEffectChannel() {
	if e.playingBit != 0 {
		e.playingBit = 0
		e.freqDir = false
		e.baseFreq = 0
		e.vol = 0
		e.command.freq = 0
	}
}

// processEffect processes the current or next effect
func (e *EffectProcessor) processEffect() byte {
	for {
		if e.queueMask == 0 {
			e.clearEffectChannel()
			break
		}

		effectId := byte(7)
		effectBit := uint8(0x80)
		for effectBit != 0 {
			if e.queueMask&effectBit != 0 {
				break
			}
			effectBit >>= 1
			effectId -= 1
		}

		e.processEffectBit(effectId, effectBit)

		if e.queueMask&effectBit != 0 {
			e.computeEffectFreq()
			e.computeEffectVol()
			break
		}
	}
	return e.vol
}

// processEffectBit starts (or continues) playing a specific effect
func (e *EffectProcessor) processEffectBit(effectId byte, effectBit uint8) {
	// processing effect yet?
	if (e.playingBit & effectBit) == 0 {
		// not yet
		e.playingBit = effectBit

		if e.index == 2 && alternateMode {
			effectId += effect2AlternateOffset
		}
		tbl := effectTable[e.index][effectId]

		// initialise state from decode table data
		e.octave = (tbl.octaveAndWave >> 4) & 0x07
		e.wave = tbl.octaveAndWave & 0x0f
		e.initialBaseFreq = tbl.initialBaseFreq
		e.freqIncr = tbl.freqIncr
		e.duration = tbl.reverseAndDuration & 0x7f
		e.reverse = tbl.reverseAndDuration&0x80 != 0
		e.repeatFreqIncr = tbl.repeatFreqIncr
		e.repeatCounter = tbl.repeatCounter
		e.initialVol = tbl.envelopeAndInitialVol & 0x0f
		e.envelope = tbl.envelopeAndInitialVol >> 4
		e.volIncr = tbl.volIncr

		e.durationCounter = e.duration
		e.baseFreq = e.initialBaseFreq

		if (e.envelope & ENV_ATTACK_BIT) == 0 {
			e.vol = e.initialVol
			e.freqDir = false
		}
	}

	// has duration been exhausted yet?
	e.durationCounter -= 1
	if e.durationCounter != 0 {
		return
	}

	// do we repeat?
	if e.repeatCounter != 0 {
		// have we finished repeating?
		e.repeatCounter -= 1
		if e.repeatCounter == 0 {
			// mark the effect as finished, so we can move onto the next, if any
			e.queueMask &= ^effectBit
			return
		}
	}

	// reset the duration
	e.durationCounter = e.duration

	// if this is a reversing effect,
	// swap the direction of frequency sweep
	if e.reverse {
		e.freqIncr = -e.freqIncr
		// recompute on odd direction changes
		e.freqDir = !e.freqDir
		if e.freqDir {
			return
		}
	}

	// here if this is an even direction change,
	// or the start of a regular repeat

	// bump the starting freq
	e.initialBaseFreq += e.repeatFreqIncr
	e.baseFreq = e.initialBaseFreq

	// bump the starting volume
	e.initialVol += e.volIncr

	// surely this is always true as e.envelope is defined as 0-4
	// In practice only 0 and 1 are used in the effects
	if (e.envelope & ENV_ATTACK_BIT) == 0 {
		e.vol = e.initialVol
	}
}

// computeEffectFreq outputs a new frequency shifted to the current octave.
func (e *EffectProcessor) computeEffectFreq() {
	e.baseFreq += e.freqIncr
	e.command.freq = uint32(e.baseFreq) << e.octave
}

// computeEffectVol modulates the output volume according to the current envelope.
func (e *EffectProcessor) computeEffectVol() {
	e.vol = applyEnvelope(e, e.vol, e.envelope)
}
