package audio

// A SongProcessor represents the processor for a single channel of a song.
type SongProcessor struct {
	counter
	index           int
	command         *Command
	queueMask       uint8
	playingBit      uint8
	envelope        byte
	durationCounter byte
	baseFreq        byte
	vol             byte
	prog            []byte
	pc              uint16
	wave            byte
	initialVol      byte
	octave          byte
	octaveUp        byte
}

// processSongs runs all the song processors.
func (au *Audio) processSongs() {
	for _, s := range au.songProcessor {
		vol := s.processSong()
		if s.queueMask != 0 {
			s.command.vol = vol
		}
	}
}

// processSong runs the processor for a single song.
func (s *SongProcessor) processSong() byte {
	if s.queueMask == 0 {
		s.clearSongChannel()
	} else {
		songId := SongId(7)
		songBit := uint8(0x80)
		for songBit != 0 {
			if s.queueMask&songBit != 0 {
				break
			}
			songBit >>= 1
			songId -= 1
		}
		s.processSongBit(songId, songBit)
	}

	return s.vol
}

// processSongBit starts (or continues) playing a specific song.
func (s *SongProcessor) processSongBit(songId SongId, songBit uint8) {

	// Have we started yet?
	if s.playingBit&songBit == 0 {

		// TODO - in alternate mode, we should behave as follows:
		// if songId == 2 {
		//   switch levelNumber {
		//     case 1:  songId = 1
		//     case 4:  songId = 2
		//     default: songId = 3
		//   }
		// } else {
		//   songId = 0
		// }

		s.playingBit = songBit
		s.prog = songTable[songId][s.index]
		s.pc = 0
	} else {
		// already playing
		s.durationCounter -= 1
		if s.durationCounter != 0 {
			s.computeSongFreq()
			s.computeSongVol()
			return
		}
	}

	s.processSongOp()
}

// processSongOp executes the next opcode or note in the current song's program.
func (s *SongProcessor) processSongOp() {
	for {
		op := s.prog[s.pc]
		s.pc++
		switch op {
		case SONG_OP_GOTO:
			lo := s.prog[s.pc]
			s.pc++
			hi := s.prog[s.pc]
			s.pc++
			s.pc = (uint16(hi) << 8) | uint16(lo)
		case SONG_OP_WAVE:
			s.wave = s.prog[s.pc]
			s.pc++
		case SONG_OP_OCTAVE:
			s.octave = s.prog[s.pc]
			s.pc++
		case SONG_OP_VOLUME:
			s.initialVol = s.prog[s.pc]
			s.pc++
		case SONG_OP_ENVELOPE:
			s.envelope = s.prog[s.pc]
			s.pc++
		case SONG_OP_END:
			s.queueMask &= ^s.playingBit
			s.clearSongChannel()
			return
		default:
			if op < SONG_OP_SPECIALS {
				s.processNote(op)
				s.computeSongFreq()
				s.computeSongVol()
				return
			}
		}
	}
}

// processNote starts playing a specific note.
func (s *SongProcessor) processNote(note byte) {
	if s.envelope&ENV_ATTACK_BIT != 0 {
		s.vol = 0
	} else {
		s.vol = s.initialVol
	}
	s.durationCounter = 1 << (note >> 5)

	// if bottom 5 bits are clear, we'll simply repeat the
	// previously played note at a given duration
	if (note & 0x1f) != 0 {
		// if the octave bit is set and the base freq bits are clear,
		// we'll play a rest, as baseFreqTable[0] == 0, and 0 freq
		// corresponds to silence.
		s.octaveUp = note & 0x10
		s.baseFreq = baseFreqTable[note&0x0f]
	}
}

// clearSongChannel stops any currently playing song.
func (s *SongProcessor) clearSongChannel() {
	if s.playingBit != 0 {
		s.playingBit = 0
		s.octaveUp = 0
		s.baseFreq = 0
		s.vol = 0
		s.command.freq = 0
	}
}

// computeSongFreq outputs a new frequency according to the current octave.
func (s *SongProcessor) computeSongFreq() {
	octave := s.octave
	if s.octaveUp != 0 {
		octave++
	}
	s.command.freq = uint32(s.baseFreq) << octave
}

// computeEffectVol modulates the output volume according to the current envelope.
func (s *SongProcessor) computeSongVol() {
	s.vol = applyEnvelope(s, s.vol, s.envelope)
}
