package audio

// Some alternate sound effects will be used when this is set
const alternateMode = false

const channelCount = 3

type Channel struct {
	freq uint32
	vol  byte
}

type Audio struct {
	songChannel   [channelCount]*SongChannel
	effectChannel [channelCount]*EffectChannel
	channel       [channelCount]Channel
	soundCounter  byte
	mute          bool

	hwVoice       [voiceCount]HwVoice
	pos           int64
	nextFrameTime float64
}

func NewAudio() *Audio {
	au := &Audio{}
	for i := range channelCount {
		channel := &au.channel[i]
		au.songChannel[i] = &SongChannel{
			index:   i,
			counter: au,
			channel: channel,
		}
		au.effectChannel[i] = &EffectChannel{
			index:   i,
			counter: au,
			channel: channel,
		}
	}
	return au
}

// implements counter interface for the benefit of the envelope generator
func (au *Audio) GetCount() byte {
	return au.soundCounter
}

func (au *Audio) Mute() {
	au.mute = true
}

func (au *Audio) UnMute() {
	au.mute = false
}

func (au *Audio) PlaySong(song int) {
	au.songChannel[0].queueMask |= 1 << song // melody
	au.songChannel[1].queueMask |= 1 << song // rhythm
}

func (au *Audio) PlayTransientEffect(i TransientEffect) {
	au.effectChannel[0].queueMask |= (1 << i)
}

func (au *Audio) StopAllTransientEffects() {
	au.effectChannel[0].queueMask = 0
}

func (au *Audio) PlayBackgroundEffect(i BackgroundEffect) {
	const backgroundMask = (1 << EnergiserEaten) | (1 << EyesReturning)
	e := au.effectChannel[1]
	background := e.queueMask & backgroundMask
	e.queueMask = (e.queueMask & background) | (1 << i)
}

func (au *Audio) StopBackgroundEffect(i BackgroundEffect) {
	au.effectChannel[1].queueMask &= ^(1 << i)
}

func (au *Audio) StopAllBackgroundEffects() {
	au.effectChannel[1].queueMask = 0
}

func (au *Audio) PlayPacmanEffect(i PacmanEffect) {
	const even = DotEatenEven
	const odd = DotEatenOdd
	const evenOddMask = byte(1)<<even | byte(1)<<odd
	e := au.effectChannel[2]
	if i == even || i == odd {
		qm := effectChannel[2].queueMask
		e.queueMask = (qm & ^evenOddMask) | (1 << i)
	} else {
		e.queueMask = (1 << i)
	}
}

func (au *Audio) StopAllPacmanEffects() {
	au.effectChannel[2].queueMask = 0
}
