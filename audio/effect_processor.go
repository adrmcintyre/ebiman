package audio

type EffectChannel struct {
	queue_mask        uint8
	playing_bit       uint8
	envelope          byte
	duration_counter  byte
	base_freq         byte
	vol               byte
	freq_dir          byte
	octave            byte
	wave              byte
	initial_base_freq byte
	freq_incr         byte
	duration          byte
	reverse           byte
	repeat_freq_incr  byte
	repeat_counter    byte
	initial_vol       byte
	vol_incr          byte
}

var effect_channel [channel_count]EffectChannel

func PlayEffect1(i int) {
	effect_channel[0].queue_mask |= (1 << i)
}

func MuteEffect1() {
	effect_channel[0].queue_mask = 0
}

func PlayEffect2(i int) {
	const backgroundMask = (1 << Effect2_EnergiserEaten) | (1 << Effect2_EyesReturning)
	background := effect_channel[1].queue_mask & backgroundMask
	effect_channel[1].queue_mask = (effect_channel[1].queue_mask & background) | (1 << i)
}

func StopEffect2(i int) {
	effect_channel[1].queue_mask &= ^(1 << i)
}

func MuteEffect2() {
	effect_channel[1].queue_mask = 0
}

func PlayEffect3(i int) {
	const even = Effect3_DotEatenEven
	const odd = Effect3_DotEatenOdd
	const evenOddMask = byte(1)<<even | byte(1)<<odd
	if i == even || i == odd {
		qm := effect_channel[2].queue_mask
		effect_channel[2].queue_mask = (qm & ^evenOddMask) | (1 << i)
	} else {
		effect_channel[2].queue_mask = (1 << i)
	}
}

func MuteEffect3() {
	effect_channel[2].queue_mask = 0
}

func process_effects() {
	for chIndex := range channel_count {
		channel[chIndex].vol = process_effect(chIndex)
	}
	channel[0].freq &= 0xffff // retain bottom 16 bits only
}

func clear_effect_channel(chIndex int) {
	e := &effect_channel[chIndex]
	if e.playing_bit != 0 {
		e.playing_bit = 0
		e.freq_dir = 0
		e.base_freq = 0
		e.vol = 0
		channel[chIndex].freq = 0
	}
}

// Process effect (one voice)
func process_effect(chIndex int) byte {
	e := &effect_channel[chIndex]
	for {
		if e.queue_mask == 0 {
			clear_effect_channel(chIndex)
			break
		}

		effect_num := byte(7)
		effect_bit := uint8(0x80)
		for effect_bit != 0 {
			if e.queue_mask&effect_bit != 0 {
				break
			}
			effect_bit >>= 1
			effect_num -= 1
		}

		process_effect_bit(chIndex, effect_num, effect_bit)

		if e.queue_mask&effect_bit != 0 {
			compute_effect_freq(chIndex)
			compute_effect_vol(chIndex)
			break
		}
	}
	return e.vol
}

// Process effect bit : process one effect, represented by 1 bit (in E)
func process_effect_bit(chIndex int, effect_num byte, effect_bit uint8) {
	e := &effect_channel[chIndex]

	// processing effect yet?
	if (e.playing_bit & effect_bit) == 0 {
		// not yet
		e.playing_bit = effect_bit

		if chIndex == 2 && audio_mspacman_mode {
			effect_num += effect2_mspacman_offset
		}
		table := effect_table[chIndex][effect_num]

		e.octave = (table[0] >> 4) & 0x07
		e.wave = table[0] & 0x0f
		e.initial_base_freq = table[1]
		e.freq_incr = table[2]
		e.duration = table[3] & 0x7f
		e.reverse = table[3] & 0x80
		e.repeat_freq_incr = table[4]
		e.repeat_counter = table[5]
		e.initial_vol = table[6] & 0x0f
		e.envelope = table[6] >> 4
		e.vol_incr = table[7]

		e.duration_counter = e.duration
		e.base_freq = e.initial_base_freq

		if (e.envelope & ENV_ATTACK_BIT) == 0 {
			e.vol = e.initial_vol
			e.freq_dir = 0
		}
	}

	// has duration been exhausted yet?
	e.duration_counter -= 1
	if e.duration_counter != 0 {
		return
	}

	// do we repeat?
	if e.repeat_counter != 0 {
		// have we finished repeating?
		e.repeat_counter -= 1
		if e.repeat_counter == 0 {
			// mark the effect as finished, so we can move onto the next, if any
			e.queue_mask &= ^effect_bit
			return
		}
	}

	// reset the duration
	e.duration_counter = e.duration

	// if this is a reversing effect,
	// swap the direction of frequency sweep
	if e.reverse != 0 {
		e.freq_incr = -e.freq_incr
		// recompute on odd direction changes
		e.freq_dir = 1 - e.freq_dir
		if e.freq_dir != 0 {
			return
		}
	}

	// here if this is an even direction change,
	// or the start of a regular repeat

	// bump the starting freq
	e.initial_base_freq += e.repeat_freq_incr
	e.base_freq = e.initial_base_freq

	// bump the starting volume
	e.initial_vol += e.vol_incr

	// surely this is always true as e.envelope is defined as 0-4
	// In practice only 0 and 1 are used in the effects
	if (e.envelope & ENV_ATTACK_BIT) == 0 {
		e.vol = e.initial_vol
	}
}

func compute_effect_freq(chIndex int) {
	e := &effect_channel[chIndex]
	e.base_freq += e.freq_incr
	channel[chIndex].freq = uint32(e.base_freq) << e.octave
}

func compute_effect_vol(chIndex int) {
	e := &effect_channel[chIndex]
	e.vol = apply_envelope(e.vol, e.envelope)
}
