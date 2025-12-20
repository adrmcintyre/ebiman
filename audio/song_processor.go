package audio

type SongChannel struct {
	queue_mask       uint8
	playing_bit      uint8
	envelope         byte
	duration_counter byte
	base_freq        byte
	vol              byte

	prog        []byte
	pc          uint16
	wave        byte
	initial_vol byte
	octave      byte
	octave_up   byte
}

var song_channel [channel_count]SongChannel

func PlaySong(song int) {
	song_channel[0].queue_mask |= 1 << song // melody
	song_channel[1].queue_mask |= 1 << song // rhythm
}

func process_songs() {
	for i := range channel_count {
		vol := process_song(i)
		if song_channel[i].queue_mask != 0 {
			channel[i].vol = vol
		}
	}
}

func process_song(chIndex int) byte {
	s := &song_channel[chIndex]
	if s.queue_mask == 0 {
		clear_song_channel(chIndex)
	} else {
		song_num := 7
		song_bit := uint8(0x80)
		for song_bit != 0 {
			if s.queue_mask&song_bit != 0 {
				break
			}
			song_bit >>= 1
			song_num -= 1
		}
		process_song_bit(chIndex, song_num, song_bit)
	}

	return s.vol
}

func process_song_bit(chIndex int, song_num int, song_bit uint8) {
	s := &song_channel[chIndex]

	// Have we started yet?
	if s.playing_bit&song_bit == 0 {
		// not started yet

		//#ifdef MSPACMAN
		//    if song_num == 2 {
		//      switch level_number {
		//        case 1:  song_num = 1
		//        case 4:  song_num = 2
		//        default: song_num = 3
		//      }
		//    } else {
		//      song_num = 0
		//    }
		//#endif

		s.playing_bit = song_bit
		s.prog = song_table[song_num][chIndex]
		s.pc = 0
	} else {
		// already playing
		s.duration_counter -= 1
		if s.duration_counter != 0 {
			compute_song_freq(chIndex)
			compute_song_vol(chIndex)
			return
		}
	}

	process_song_op(chIndex)
}

func process_song_op(chIndex int) {
	s := &song_channel[chIndex]

	for {
		op := s.prog[s.pc]
		s.pc++
		if op < SONG_OP_SPECIALS {
			process_regular_op(chIndex, op)
			compute_song_freq(chIndex)
			compute_song_vol(chIndex)
			return
		} else if op == SONG_OP_REPEAT {
			lo := s.prog[s.pc]
			s.pc++
			hi := s.prog[s.pc]
			s.pc++
			s.pc = (uint16(hi) << 8) | uint16(lo)
		} else if op == SONG_OP_WAVE {
			s.wave = s.prog[s.pc]
			s.pc++
		} else if op == SONG_OP_OCTAVE {
			s.octave = s.prog[s.pc]
			s.pc++
		} else if op == SONG_OP_VOLUME {
			s.initial_vol = s.prog[s.pc]
			s.pc++
		} else if op == SONG_OP_ENVELOPE {
			s.envelope = s.prog[s.pc]
			s.pc++
		} else if op == SONG_OP_END {
			s.queue_mask &= ^s.playing_bit
			clear_song_channel(chIndex)
			return
		} else {
			// 0xf5 .. 0xfe : nop
		}
	}
}

func process_regular_op(chIndex int, op byte) {
	s := &song_channel[chIndex]

	if s.envelope&ENV_ATTACK_BIT != 0 {
		s.vol = 0
	} else {
		s.vol = s.initial_vol
	}
	s.duration_counter = 1 << (op >> 5)

	// if bottom 5 bits are clear, we'll simply repeated the
	// previously played note at a given duration
	if (op & 0x1f) != 0 {
		// if the octave bit is set and the base freq bits are clear,
		// we'll play a rest, as base_freq_table[0] == 0, and 0 freq
		// corresponds to silence.
		s.octave_up = op & 0x10
		s.base_freq = baseFreqTable[op&0x0f]
	}
}

func clear_song_channel(chIndex int) {
	s := &song_channel[chIndex]
	if s.playing_bit != 0 {
		s.playing_bit = 0
		s.octave_up = 0
		s.base_freq = 0
		s.vol = 0
		channel[chIndex].freq = 0
	}
}

func compute_song_freq(chIndex int) {
	s := &song_channel[chIndex]
	octave := s.octave
	if s.octave_up != 0 {
		octave++
	}
	channel[chIndex].freq = uint32(s.base_freq) << octave
}

func compute_song_vol(chIndex int) {
	s := &song_channel[chIndex]
	s.vol = apply_envelope(s.vol, s.envelope)
}
