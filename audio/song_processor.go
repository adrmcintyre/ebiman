package audio

type SongChannel struct {
	counter
	index           int
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
	channel         *Channel
}

func (au *Audio) processSongs() {
	for _, s := range au.songChannel {
		vol := s.processSong()
		if s.queueMask != 0 {
			s.channel.vol = vol
		}
	}
}

func (s *SongChannel) processSong() byte {
	if s.queueMask == 0 {
		s.clearSongChannel()
	} else {
		songNum := 7
		songBit := uint8(0x80)
		for songBit != 0 {
			if s.queueMask&songBit != 0 {
				break
			}
			songBit >>= 1
			songNum -= 1
		}
		s.processSongBit(songNum, songBit)
	}

	return s.vol
}

func (s *SongChannel) processSongBit(songNum int, songBit uint8) {

	// Have we started yet?
	if s.playingBit&songBit == 0 {

		// TODO - in alternate mode, we should behave as follows:
		// if songNum == 2 {
		//   switch levelNumber {
		//     case 1:  songNum = 1
		//     case 4:  songNum = 2
		//     default: songNum = 3
		//   }
		// } else {
		//   songNum = 0
		// }

		s.playingBit = songBit
		s.prog = songTable[songNum][s.index]
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

func (s *SongChannel) processSongOp() {
	for {
		op := s.prog[s.pc]
		s.pc++
		if op < SONG_OP_SPECIALS {
			s.processRegularOp(op)
			s.computeSongFreq()
			s.computeSongVol()
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
			s.initialVol = s.prog[s.pc]
			s.pc++
		} else if op == SONG_OP_ENVELOPE {
			s.envelope = s.prog[s.pc]
			s.pc++
		} else if op == SONG_OP_END {
			s.queueMask &= ^s.playingBit
			s.clearSongChannel()
			return
		} else {
			// 0xf5 .. 0xfe : nop
		}
	}
}

func (s *SongChannel) processRegularOp(op byte) {
	if s.envelope&ENV_ATTACK_BIT != 0 {
		s.vol = 0
	} else {
		s.vol = s.initialVol
	}
	s.durationCounter = 1 << (op >> 5)

	// if bottom 5 bits are clear, we'll simply repeated the
	// previously played note at a given duration
	if (op & 0x1f) != 0 {
		// if the octave bit is set and the base freq bits are clear,
		// we'll play a rest, as baseFreqTable[0] == 0, and 0 freq
		// corresponds to silence.
		s.octaveUp = op & 0x10
		s.baseFreq = baseFreqTable[op&0x0f]
	}
}

func (s *SongChannel) clearSongChannel() {
	if s.playingBit != 0 {
		s.playingBit = 0
		s.octaveUp = 0
		s.baseFreq = 0
		s.vol = 0
		s.channel.freq = 0
	}
}

func (s *SongChannel) computeSongFreq() {
	octave := s.octave
	if s.octaveUp != 0 {
		octave++
	}
	s.channel.freq = uint32(s.baseFreq) << octave
}

func (s *SongChannel) computeSongVol() {
	s.vol = applyEnvelope(s, s.vol, s.envelope)
}
