package audio

// table of note frequencies
var baseFreqTable = [16]byte{
	0x00, 0x57, 0x5c, 0x61, 0x67, 0x6d, 0x74, 0x7b,
	0x82, 0x8a, 0x92, 0x9a, 0xa3, 0xad, 0xb8, 0xc3,
}

var hw_voice [3]HwVoice

// Some alternate sound effects will be used when this is set
var audio_mspacman_mode = false

type Channel struct {
	freq uint32
	vol  byte
}

const channel_count = 3

var sound_counter byte
var channel [channel_count]Channel

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
