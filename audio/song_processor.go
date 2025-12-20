package audio

type SongChannel struct {
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

var songChannel [channelCount]SongChannel

func PlaySong(song int) {
	songChannel[0].queueMask |= 1 << song // melody
	songChannel[1].queueMask |= 1 << song // rhythm
}

func processSongs() {
	for i := range channelCount {
		vol := processSong(i)
		if songChannel[i].queueMask != 0 {
			channel[i].vol = vol
		}
	}
}

func processSong(chIndex int) byte {
	s := &songChannel[chIndex]
	if s.queueMask == 0 {
		clearSongChannel(chIndex)
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
		processSongBit(chIndex, songNum, songBit)
	}

	return s.vol
}

func processSongBit(chIndex int, songNum int, songBit uint8) {
	s := &songChannel[chIndex]

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
		s.prog = songTable[songNum][chIndex]
		s.pc = 0
	} else {
		// already playing
		s.durationCounter -= 1
		if s.durationCounter != 0 {
			computeSongFreq(chIndex)
			computeSongVol(chIndex)
			return
		}
	}

	processSongOp(chIndex)
}

func processSongOp(chIndex int) {
	s := &songChannel[chIndex]

	for {
		op := s.prog[s.pc]
		s.pc++
		if op < SONG_OP_SPECIALS {
			processRegularOp(chIndex, op)
			computeSongFreq(chIndex)
			computeSongVol(chIndex)
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
			clearSongChannel(chIndex)
			return
		} else {
			// 0xf5 .. 0xfe : nop
		}
	}
}

func processRegularOp(chIndex int, op byte) {
	s := &songChannel[chIndex]

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

func clearSongChannel(chIndex int) {
	s := &songChannel[chIndex]
	if s.playingBit != 0 {
		s.playingBit = 0
		s.octaveUp = 0
		s.baseFreq = 0
		s.vol = 0
		channel[chIndex].freq = 0
	}
}

func computeSongFreq(chIndex int) {
	s := &songChannel[chIndex]
	octave := s.octave
	if s.octaveUp != 0 {
		octave++
	}
	channel[chIndex].freq = uint32(s.baseFreq) << octave
}

func computeSongVol(chIndex int) {
	s := &songChannel[chIndex]
	s.vol = applyEnvelope(s.vol, s.envelope)
}
