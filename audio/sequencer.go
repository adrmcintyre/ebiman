package audio

// nudgeSequencer advances the state of the sound processors,
// and updates the hardware audio output.
// It should be called at 60Hz.
func (au *Audio) nudgeSequencer() {
	au.writeRegisters()

	// advance sound state
	au.soundCounter += 1
	au.processEffects()
	au.processSongs()
}

// writeRegisters updates the simulate hardware registers
// according to the active songs and effects - not that
// song output takes precedence over any active effects.
func (au *Audio) writeRegisters() {
	for chIndex := range channelCount {
		if au.songProcessor[chIndex].queueMask != 0 {
			au.hwVoice[chIndex].wave = au.songProcessor[chIndex].wave
		} else {
			au.hwVoice[chIndex].wave = au.effectProcessor[chIndex].wave
		}
		au.hwVoice[chIndex].freq = au.command[chIndex].freq
		au.hwVoice[chIndex].vol = au.command[chIndex].vol
	}
}
