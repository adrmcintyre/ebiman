package audio

const (
	MIN_OUTPUT_VOLUME     = 0
	MAX_OUTPUT_VOLUME     = 5
	DEFAULT_OUTPUT_VOLUME = 1 // quietest non-zero volume
)

// outputVolumes are volume settings as a proportion
// of the host's maximum volume level.
var outputVolumes = [MAX_OUTPUT_VOLUME + 1]float64{
	0.000,
	0.025,
	0.050,
	0.075,
	0.100,
	0.125,
}

func (au *Audio) SetOutputVolume(volume int) {
	au.outputVolume = max(MIN_OUTPUT_VOLUME, min(volume, MAX_OUTPUT_VOLUME))
	au.player.SetVolume(outputVolumes[au.outputVolume])
}

func (au *Audio) OutputVolumeUp() {
	au.SetOutputVolume(au.outputVolume + 1)
}

func (au *Audio) OutputVolumeDown() {
	au.SetOutputVolume(au.outputVolume - 1)
}

// Mute mutes all audio output.
func (au *Audio) Mute() {
	au.muted = true
}

// UnMute un-mutes the audio output.
func (au *Audio) UnMute() {
	au.muted = false
}
