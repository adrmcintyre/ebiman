package audio

type EffectChannel struct {
	counter
	index           int
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
	channel         *Channel
}

var effectChannel [channelCount]EffectChannel

func (au *Audio) processAllEffects() {
	for _, e := range au.effectChannel {
		e.channel.vol = e.processEffect()
	}
	au.channel[0].freq &= 0xffff // retain bottom 16 bits only
}

func (e *EffectChannel) clearEffectChannel() {
	if e.playingBit != 0 {
		e.playingBit = 0
		e.freqDir = false
		e.baseFreq = 0
		e.vol = 0
		e.channel.freq = 0
	}
}

// Process effect (one voice)
func (e *EffectChannel) processEffect() byte {
	for {
		if e.queueMask == 0 {
			e.clearEffectChannel()
			break
		}

		effectNum := byte(7)
		effectBit := uint8(0x80)
		for effectBit != 0 {
			if e.queueMask&effectBit != 0 {
				break
			}
			effectBit >>= 1
			effectNum -= 1
		}

		e.processEffectBit(effectNum, effectBit)

		if e.queueMask&effectBit != 0 {
			e.computeEffectFreq()
			e.computeEffectVol()
			break
		}
	}
	return e.vol
}

// Process effect bit : process one effect, represented by 1 bit (in E)
func (e *EffectChannel) processEffectBit(effectNum byte, effectBit uint8) {
	// processing effect yet?
	if (e.playingBit & effectBit) == 0 {
		// not yet
		e.playingBit = effectBit

		if e.index == 2 && alternateMode {
			effectNum += effect2AlternateOffset
		}
		tbl := effectTable[e.index][effectNum]

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

func (e *EffectChannel) computeEffectFreq() {
	e.baseFreq += e.freqIncr
	e.channel.freq = uint32(e.baseFreq) << e.octave
}

func (e *EffectChannel) computeEffectVol() {
	e.vol = applyEnvelope(e, e.vol, e.envelope)
}
