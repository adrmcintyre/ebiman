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

// processSongs runs all the song processors.
func (au *Audio) processSongs() {
	for _, s := range au.songProcessor {
		s.processSong()
	}
}

// processEffects runs all the effects processors
func (au *Audio) processEffects() {
	for _, e := range au.effectProcessor {
		e.processEffect()
	}
}

// writeRegisters updates the simulate hardware registers
// according to the active songs and effects - not that
// song output takes precedence over any active effects.
func (au *Audio) writeRegisters() {
	au.command[0].freq &= 0xffff // retain bottom 16 bits only
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
