package audio

const (
	MinOutputVolume     = 0
	MaxOutputVolume     = 5
	DefaultOutputVolume = 3
)

// outputVolumes are volume settings as a proportion
// of the host's maximum volume level.
var outputVolumes = [MaxOutputVolume + 1]float64{
	0.00,
	0.06,
	0.12,
	0.25,
	0.50,
	1.00,
}

func (au *Audio) SetOutputVolume(volume int) {
	au.outputVolume = max(MinOutputVolume, min(volume, MaxOutputVolume))
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
