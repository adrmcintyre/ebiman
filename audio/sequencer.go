package audio

// 0 : b6-4: frequency shift | b2-0: wave select
// 1 : initial base frequency
// 2 : frequency increment (added to base freq)
// 3 : b7: reverse | b6-0: duration
// 4 : frequency increment (added to initial base frequency). Used when repeat > 1
// 5 : repeat
// 6 : b7-4: volume adjust type | b3-0: initial volume
// 7 : volume increment
type effectData [8]byte

const (
	Effect1_ExtraLife = iota
	Effect1_Credit
)

var effects_table1 = []effectData{
	{0x73, 0x20, 0x00, 0x0c, 0x00, 0x0a, 0x1f, 0x00}, // 0 - extra life
	{0x72, 0x20, 0xfb, 0x87, 0x00, 0x02, 0x0f, 0x00}, // 1 - credit
}

const (
	Effect2_EndEnergiser = iota
	Effect2_1
	Effect2_2
	Effect2_3
	Effect2_4
	Effect2_EnergiserEaten
	Effect2_EyesReturning
	Effect2_Unused
)

const effect2_mspacman_offset = 8

var effects_table2 = []effectData{
	//pacman
	{0x36, 0x20, 0x04, 0x8c, 0x00, 0x00, 0x06, 0x00}, // 0 - end of energizer
	{0x36, 0x28, 0x05, 0x8b, 0x00, 0x00, 0x06, 0x00}, // 1 - higher freq when 155 dots eaten
	{0x36, 0x30, 0x06, 0x8a, 0x00, 0x00, 0x06, 0x00}, // 2 - higher freq when 179 dots eaten
	{0x36, 0x3c, 0x07, 0x89, 0x00, 0x00, 0x06, 0x00}, // 3 - higher freq when 12 dots left
	{0x36, 0x48, 0x08, 0x88, 0x00, 0x00, 0x06, 0x00}, // 4 - reset higher freq when 12 or less dots left
	{0x24, 0x00, 0x06, 0x08, 0x00, 0x00, 0x0a, 0x00}, // 5 - energizer eaten
	{0x40, 0x70, 0xfa, 0x10, 0x00, 0x00, 0x0a, 0x00}, // 6 - eyes returning
	{0x70, 0x04, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00}, // 7 - unused

	//mspacman
	{0x59, 0x01, 0x06, 0x08, 0x00, 0x00, 0x02, 0x00}, // 0 - end of energizer
	{0x59, 0x01, 0x06, 0x09, 0x00, 0x00, 0x02, 0x00}, // 1 - higher freq when 155 dots eaten
	{0x59, 0x02, 0x06, 0x0a, 0x00, 0x00, 0x02, 0x00}, // 2 - higher freq when 179 dots eaten
	{0x59, 0x03, 0x06, 0x0b, 0x00, 0x00, 0x02, 0x00}, // 3 - higher freq when 12 dots left
	{0x59, 0x04, 0x06, 0x0c, 0x00, 0x06, 0x02, 0x00}, // 4 - reset higher freq when 12 or less dots left
	{0x24, 0x00, 0x06, 0x08, 0x02, 0x00, 0x0a, 0x00}, // 5 - energizer eaten
	{0x36, 0x07, 0x87, 0x6f, 0x00, 0x00, 0x04, 0x00}, // 6 - eyes returning
	{0x70, 0x04, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00}, // 7 - unused
}

const (
	Effect3_DotEatenEven = iota
	Effect3_DotEatenOdd
	Effect3_FruitEaten
	Effect3_GhostEaten
	Effect3_PacmanDead
	Effect3_PacmanPop
)

// pacman
var effects_table3 = []effectData{
	{0x42, 0x18, 0xfd, 0x06, 0x00, 0x01, 0x0c, 0x00}, // 0 - dot eating sound 1
	{0x42, 0x04, 0x03, 0x06, 0x00, 0x01, 0x0c, 0x00}, // 1 - dot eating sound 2
	{0x56, 0x0c, 0xff, 0x8c, 0x00, 0x02, 0x0f, 0x00}, // 2 - fruit eating sound
	{0x05, 0x00, 0x02, 0x20, 0x00, 0x01, 0x0c, 0x00}, // 3 - blue ghost eaten
	{0x41, 0x20, 0xff, 0x86, 0xfe, 0x1c, 0x0f, 0xff}, // 4 - pacman dying; ghosts bumping in act 2
	{0x70, 0x00, 0x01, 0x0c, 0x00, 0x01, 0x08, 0x00}, // 5 - pacman "pop"
}

// mspacman
var effects_table4 = []effectData{
	{0x1c, 0x70, 0x8b, 0x08, 0x00, 0x01, 0x06, 0x00}, // 0 - dot eating sound 1
	{0x1c, 0x70, 0x8b, 0x08, 0x00, 0x01, 0x06, 0x00}, // 1 - dot eating sound 2
	{0x56, 0x0c, 0xff, 0x8c, 0x00, 0x02, 0x08, 0x00}, // 2 - fruit eating sound
	{0x56, 0x00, 0x02, 0x0a, 0x07, 0x03, 0x0c, 0x00}, // 3 - blue ghost eaten
	{0x36, 0x38, 0xfe, 0x12, 0xf8, 0x04, 0x0f, 0xfc}, // 4 - pacman dying; ghosts bumping in act 2
	{0x22, 0x01, 0x01, 0x06, 0x00, 0x01, 0x07, 0x00}, // 5 - fruit bouncing
}

var effect_table = [4][]effectData{
	effects_table1,
	effects_table2,
	effects_table3,
	effects_table4,
}

var powerOf2 = [8]byte{
	0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80,
}

var base_freq_table = [16]byte{
	0x00, 0x57, 0x5c, 0x61, 0x67, 0x6d, 0x74, 0x7b,
	0x82, 0x8a, 0x92, 0x9a, 0xa3, 0xad, 0xb8, 0xc3,
}

const (
	SONG_OP_SPECIALS byte = 0xf0
	SONG_OP_REPEAT   byte = 0xf0
	SONG_OP_WAVE     byte = 0xf1
	SONG_OP_OCTAVE   byte = 0xf2
	SONG_OP_VOLUME   byte = 0xf3
	SONG_OP_ENVELOPE byte = 0xf4
	SONG_OP_END      byte = 0xf5
)

const (
	ENV_CONST      byte = 0
	ENV_DECAY1     byte = 1
	ENV_DECAY2     byte = 2
	ENV_DECAY4     byte = 3
	ENV_DECAY8     byte = 4
	ENV_ATTACK_BIT byte = 8 // partial implementation
)

var song_startup_melody = []byte{
	SONG_OP_WAVE, 0x02,
	SONG_OP_OCTAVE, 0x03,
	SONG_OP_VOLUME, 0x0f,
	SONG_OP_ENVELOPE, ENV_DECAY1,

	0x82, 0x70, 0x69,
	0x82, 0x70, 0x69,
	0x83, 0x70, 0x6a,
	0x83, 0x70, 0x6a,

	0x82, 0x70, 0x69,
	0x82, 0x70, 0x69,
	0x89, 0x8b,
	0x8d, 0x8e,

	SONG_OP_END,
}

var song_startup_rhythm = []byte{
	SONG_OP_WAVE, 0x00,
	SONG_OP_OCTAVE, 0x02,
	SONG_OP_VOLUME, 0x0f,
	SONG_OP_ENVELOPE, ENV_CONST,

	0x42, 0x50, 0x4e, 0x50,
	0x49, 0x50, 0x46, 0x50,
	0x4e, 0x49, 0x70,
	0x66, 0x70,

	0x43, 0x50, 0x4f, 0x50,
	0x4a, 0x50, 0x47, 0x50,
	0x4f, 0x4a, 0x70,
	0x67, 0x70,

	0x42, 0x50, 0x4e, 0x50,
	0x49, 0x50, 0x46, 0x50,
	0x4e, 0x49, 0x70,
	0x66, 0x70,

	0x45, 0x46, 0x47, 0x50,
	0x47, 0x48, 0x49, 0x50,
	0x49, 0x4a, 0x4b, 0x50,
	0x6e,

	SONG_OP_END,
}

var song_unused = []byte{
	SONG_OP_END,
}

var song_intermission_melody = []byte{
	SONG_OP_WAVE, 0x02,
	SONG_OP_OCTAVE, 0x03,
	SONG_OP_VOLUME, 0x0f,
	SONG_OP_ENVELOPE, ENV_DECAY1,

	0x67, 0x50, 0x30, 0x47, 0x30,
	0x67, 0x50, 0x30, 0x47, 0x30,
	0x67, 0x50, 0x30, 0x47, 0x30,
	0x4b, 0x10, 0x4c, 0x10, 0x4d, 0x10, 0x4e, 0x10,

	0x67, 0x50, 0x30, 0x47, 0x30,
	0x67, 0x50, 0x30, 0x47, 0x30,
	0x67, 0x50, 0x30, 0x47, 0x30,
	0x4b, 0x10, 0x4c, 0x10, 0x4d, 0x10, 0x4e, 0x10,

	0x67, 0x50, 0x30, 0x47, 0x30,
	0x67, 0x50, 0x30, 0x47, 0x30,
	0x67, 0x50, 0x30, 0x47, 0x30,
	0x4b, 0x10, 0x4c, 0x10, 0x4d, 0x10, 0x4e, 0x10,

	0x77, 0x20, 0x4e, 0x10, 0x4d, 0x10,
	0x4c, 0x10, 0x4a, 0x10, 0x47, 0x10, 0x46, 0x10,
	0x65, 0x30, 0x66, 0x30,
	0x67, 0x40, 0x70,

	SONG_OP_REPEAT, 0x08, 0x00,
}

var song_intermission_rhythm = []byte{
	SONG_OP_WAVE, 0x01,
	SONG_OP_OCTAVE, 0x01,
	SONG_OP_VOLUME, 0x0f,
	SONG_OP_ENVELOPE, ENV_CONST,

	0x26, 0x67, 0x26, 0x67, 0x26, 0x67, 0x23, 0x44,
	0x42, 0x47, 0x30, 0x67, 0x2a, 0x8b, 0x70, 0x26,
	0x67, 0x26, 0x67, 0x26, 0x67, 0x23, 0x44, 0x42,
	0x47, 0x30, 0x67, 0x23, 0x84, 0x70, 0x26, 0x67,
	0x26, 0x67, 0x26, 0x67, 0x23, 0x44, 0x42, 0x47,
	0x30, 0x67, 0x29, 0x6a, 0x2b, 0x6c, 0x30, 0x2c,
	0x6d, 0x40, 0x2b, 0x6c, 0x29, 0x6a, 0x67, 0x20,
	0x29, 0x6a, 0x40, 0x26, 0x87, 0x70,

	SONG_OP_REPEAT, 0x08, 0x00,
}

// mspacman song data
var song_mspacman_startup_melody = []byte{
	SONG_OP_WAVE, 0x00,
	SONG_OP_OCTAVE, 0x02,
	SONG_OP_VOLUME, 0x0a,
	SONG_OP_ENVELOPE, ENV_CONST,
	0x41, 0x43, 0x45, 0x86, 0x8a, 0x88, 0x8b, 0x6a,
	0x6b, 0x71, 0x6a, 0x88, 0x8b, 0x6a, 0x6b, 0x71,
	0x6a, 0x6b, 0x71, 0x73, 0x75, 0x96, 0x95, 0x96,
	SONG_OP_END,
}

var song_mspacman_startup_rhythm = []byte{
	SONG_OP_WAVE, 0x02,
	SONG_OP_OCTAVE, 0x03,
	SONG_OP_VOLUME, 0x0a,
	SONG_OP_ENVELOPE, ENV_DECAY2,
	0x50, 0x70, 0x86, 0x90, 0x81, 0x90, 0x86, 0x90,
	0x68, 0x6a, 0x6b, 0x68, 0x6a, 0x68, 0x66, 0x6a,
	0x68, 0x66, 0x65, 0x68, 0x86, 0x81, 0x86,
	SONG_OP_END,
}

var song_mspacman_act1_melody = []byte{
	SONG_OP_WAVE, 0x00,
	SONG_OP_OCTAVE, 0x02,
	SONG_OP_VOLUME, 0x0a,
	SONG_OP_ENVELOPE, ENV_CONST,
	0x69, 0x6b, 0x69, 0x86, 0x61, 0x64, 0x65, 0x86,
	0x86, 0x64, 0x66, 0x64, 0x61, 0x69, 0x6b, 0x69,
	0x86, 0x61, 0x64, 0x64, 0xa1, 0x70, 0x71, 0x74,
	0x75, 0x35, 0x76, 0x30, 0x50, 0x35, 0x76, 0x30,
	0x50, 0x54, 0x56, 0x54, 0x51, 0x6b, 0x69, 0x6b,
	0x69, 0x6b, 0x91, 0x6b, 0x69, 0x66, 0xf2, 0x01,
	0x74, 0x76, 0x74, 0x71, 0x74, 0x71, 0x6b, 0x69,
	0xa6, 0xa6,
	SONG_OP_END,
}

var song_mspacman_act1_rhythm = []byte{
	SONG_OP_WAVE, 0x03,
	SONG_OP_OCTAVE, 0x03,
	SONG_OP_VOLUME, 0x0a,
	SONG_OP_ENVELOPE, ENV_DECAY2,
	0x70, 0x66, 0x70, 0x46, 0x50, 0x86, 0x90, 0x70,
	0x66, 0x70, 0x46, 0x50, 0x86, 0x90, 0x70, 0x66,
	0x70, 0x46, 0x50, 0x86, 0x90, 0x70, 0x61, 0x70,
	0x41, 0x50, 0x81, 0x90, 0xf4, 0x00, 0xa6, 0xa4,
	0xa2, 0xa1, 0xf4, 0x01, 0x86, 0x89, 0x8b, 0x81,
	0x74, 0x71, 0x6b, 0x69, 0xa6,
	SONG_OP_END,
}

var song_mspacman_act3_melody = []byte{
	SONG_OP_WAVE, 0x00,
	SONG_OP_OCTAVE, 0x02,
	SONG_OP_VOLUME, 0x0a,
	SONG_OP_ENVELOPE, ENV_CONST,
	0x65, 0x64, 0x65, 0x88, 0x67, 0x88, 0x61, 0x63,
	0x64, 0x85, 0x64, 0x85, 0x6a, 0x69, 0x6a, 0x8c,
	0x75, 0x93, 0x90, 0x91, 0x90, 0x91, 0x70, 0x8a,
	0x68, 0x71,
	SONG_OP_END,
}

var song_mspacman_act3_rhythm = []byte{
	SONG_OP_WAVE, 0x02,
	SONG_OP_OCTAVE, 0x03,
	SONG_OP_VOLUME, 0x0a,
	SONG_OP_ENVELOPE, ENV_DECAY2,
	0x65, 0x90, 0x68, 0x70, 0x68, 0x67, 0x66, 0x65,
	0x90, 0x61, 0x70, 0x61, 0x65, 0x68, 0x66, 0x90,
	0x63, 0x90, 0x86, 0x90, 0x85, 0x90, 0x85, 0x70,
	0x86, 0x68, 0x65,
	SONG_OP_END,
}

const (
	SongStartup = iota
	SongIntermission
	SongMspacmanStartup
	SongMspacmanAct1
	SongMspacmanAct2
	SongMspacmanAct3
)

var song_table = [6][4][]byte{
	// pacman
	{
		song_startup_melody,
		song_startup_rhythm,
		song_unused,
		song_unused,
	},
	{
		song_intermission_melody,
		song_intermission_rhythm,
		song_unused,
		song_unused,
	},
	// mspacman
	{
		song_mspacman_startup_melody,
		song_mspacman_startup_rhythm,
		song_unused,
		song_unused,
	},
	{
		song_mspacman_act1_melody,
		song_mspacman_act1_rhythm,
		song_unused,
		song_unused,
	},
	{
		song_intermission_melody,
		song_intermission_rhythm,
		song_unused,
		song_unused,
	},
	{
		song_mspacman_act3_melody,
		song_mspacman_act3_rhythm,
		song_unused,
		song_unused,
	},
}

var hw_voice [3]HwVoice

// Some alternate sound effects will be used when this is set
var audio_mspacman_mode = false

type Channel struct {
	freq uint32
	vol  byte
}

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

const channel_count = 3

var sound_counter byte
var channel [channel_count]Channel
var song_channel [channel_count]SongChannel
var effect_channel [channel_count]EffectChannel

func PlaySong(song int) {
	song_channel[0].queue_mask |= 1 << song // melody
	song_channel[1].queue_mask |= 1 << song // rhythm
}

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
		s.base_freq = base_freq_table[op&0x0f]
	}
}

func process_effects() {
	for chIndex := range channel_count {
		channel[chIndex].vol = process_effect(chIndex)
	}
	channel[0].freq &= 0xffff // retain bottom 16 bits only
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

func compute_song_freq(chIndex int) {
	s := &song_channel[chIndex]
	octave := s.octave
	if s.octave_up != 0 {
		octave++
	}
	channel[chIndex].freq = uint32(s.base_freq) << octave
}

func compute_effect_freq(chIndex int) {
	e := &effect_channel[chIndex]
	e.base_freq += e.freq_incr
	channel[chIndex].freq = uint32(e.base_freq) << e.octave
}

func compute_song_vol(chIndex int) {
	s := &song_channel[chIndex]
	s.vol = compute_vol(s.vol, s.envelope)
}

func compute_effect_vol(chIndex int) {
	e := &effect_channel[chIndex]
	e.vol = compute_vol(e.vol, e.envelope)
}

func compute_vol(vol byte, envelope byte) byte {
	switch envelope {
	case ENV_CONST:
		return vol

	case ENV_DECAY1, ENV_DECAY2, ENV_DECAY4, ENV_DECAY8:
		rate := byte(1) << (envelope - 1)
		if sound_counter&(rate-1) == 0 {
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

func Setup() {
	for chIndex := range channel_count {
		hw_voice[chIndex].wave = 0
		hw_voice[chIndex].vol = 0
		hw_voice[chIndex].freq = 0
		song_channel[chIndex].queue_mask = 0
		effect_channel[chIndex].queue_mask = 0
	}
}

// load the sound into the hardware
func audio_hw_write() {
	for chIndex := range channel_count {
		if song_channel[chIndex].queue_mask != 0 {
			hw_voice[chIndex].wave = song_channel[chIndex].wave
		} else {
			hw_voice[chIndex].wave = effect_channel[chIndex].wave
		}
		hw_voice[chIndex].freq = channel[chIndex].freq
		hw_voice[chIndex].vol = channel[chIndex].vol
	}
}

func audio_vblank1() {
	audio_hw_write()
}

// advance sound state
func audio_vblank2() {
	sound_counter += 1
	process_effects()
	process_songs()
}

// Should be called at 60Hz
func runSequencerFrame() {
	audio_vblank1()
	audio_vblank2()
}
