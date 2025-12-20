package audio

type EffectChannel struct {
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

var effectChannel [channelCount]EffectChannel

func PlayTransientEffect(i TransientEffect) {
	effectChannel[0].queueMask |= (1 << i)
}

func StopAllTransientEffects() {
	effectChannel[0].queueMask = 0
}

func PlayBackgroundEffect(i BackgroundEffect) {
	const backgroundMask = (1 << EnergiserEaten) | (1 << EyesReturning)
	background := effectChannel[1].queueMask & backgroundMask
	effectChannel[1].queueMask = (effectChannel[1].queueMask & background) | (1 << i)
}

func StopBackgroundEffect(i BackgroundEffect) {
	effectChannel[1].queueMask &= ^(1 << i)
}

func StopAllBackgroundEffects() {
	effectChannel[1].queueMask = 0
}

func PlayPacmanEffect(i PacmanEffect) {
	const even = DotEatenEven
	const odd = DotEatenOdd
	const evenOddMask = byte(1)<<even | byte(1)<<odd
	if i == even || i == odd {
		qm := effectChannel[2].queueMask
		effectChannel[2].queueMask = (qm & ^evenOddMask) | (1 << i)
	} else {
		effectChannel[2].queueMask = (1 << i)
	}
}

func StopAllPacmanEffects() {
	effectChannel[2].queueMask = 0
}

func processAllEffects() {
	for chIndex := range channelCount {
		channel[chIndex].vol = processEffect(chIndex)
	}
	channel[0].freq &= 0xffff // retain bottom 16 bits only
}

func clearEffectChannel(chIndex int) {
	e := &effectChannel[chIndex]
	if e.playingBit != 0 {
		e.playingBit = 0
		e.freqDir = false
		e.baseFreq = 0
		e.vol = 0
		channel[chIndex].freq = 0
	}
}

// Process effect (one voice)
func processEffect(chIndex int) byte {
	e := &effectChannel[chIndex]
	for {
		if e.queueMask == 0 {
			clearEffectChannel(chIndex)
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

		processEffectBit(chIndex, effectNum, effectBit)

		if e.queueMask&effectBit != 0 {
			computeEffectFreq(chIndex)
			computeEffectVol(chIndex)
			break
		}
	}
	return e.vol
}

// Process effect bit : process one effect, represented by 1 bit (in E)
func processEffectBit(chIndex int, effectNum byte, effectBit uint8) {
	e := &effectChannel[chIndex]

	// processing effect yet?
	if (e.playingBit & effectBit) == 0 {
		// not yet
		e.playingBit = effectBit

		if chIndex == 2 && alternateMode {
			effectNum += effect2AlternateOffset
		}
		table := effectTable[chIndex][effectNum]

		e.octave = (table[0] >> 4) & 0x07
		e.wave = table[0] & 0x0f
		e.initialBaseFreq = table[1]
		e.freqIncr = table[2]
		e.duration = table[3] & 0x7f
		e.reverse = table[3]&0x80 != 0
		e.repeatFreqIncr = table[4]
		e.repeatCounter = table[5]
		e.initialVol = table[6] & 0x0f
		e.envelope = table[6] >> 4
		e.volIncr = table[7]

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

func computeEffectFreq(chIndex int) {
	e := &effectChannel[chIndex]
	e.baseFreq += e.freqIncr
	channel[chIndex].freq = uint32(e.baseFreq) << e.octave
}

func computeEffectVol(chIndex int) {
	e := &effectChannel[chIndex]
	e.vol = applyEnvelope(e.vol, e.envelope)
}
