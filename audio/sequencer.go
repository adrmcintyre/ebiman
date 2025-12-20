package audio

// Should be called at 60Hz
func (au *Audio) nudgeSequencer() {
	au.writeRegisters()

	// advance sound state
	au.soundCounter += 1
	au.processAllEffects()
	au.processSongs()
}

// load the sound into the hardware
func (au *Audio) writeRegisters() {
	for chIndex := range channelCount {
		if au.songChannel[chIndex].queueMask != 0 {
			au.hwVoice[chIndex].wave = au.songChannel[chIndex].wave
		} else {
			au.hwVoice[chIndex].wave = au.effectChannel[chIndex].wave
		}
		au.hwVoice[chIndex].freq = au.channel[chIndex].freq
		au.hwVoice[chIndex].vol = au.channel[chIndex].vol
	}
}
