package audio

var hwVoice [3]HwVoice

// Some alternate sound effects will be used when this is set
var alternateMode = false

type Channel struct {
	freq uint32
	vol  byte
}

const channelCount = 3

var soundCounter byte
var channel [channelCount]Channel

func Setup() {
	for chIndex := range channelCount {
		hwVoice[chIndex].wave = 0
		hwVoice[chIndex].vol = 0
		hwVoice[chIndex].freq = 0
		songChannel[chIndex].queueMask = 0
		effectChannel[chIndex].queueMask = 0
	}
}

// load the sound into the hardware
func writeRegisters() {
	for chIndex := range channelCount {
		if songChannel[chIndex].queueMask != 0 {
			hwVoice[chIndex].wave = songChannel[chIndex].wave
		} else {
			hwVoice[chIndex].wave = effectChannel[chIndex].wave
		}
		hwVoice[chIndex].freq = channel[chIndex].freq
		hwVoice[chIndex].vol = channel[chIndex].vol
	}
}

// Should be called at 60Hz
func nudgeSequencer() {
	writeRegisters()

	// advance sound state
	soundCounter += 1
	processAllEffects()
	processSongs()
}
