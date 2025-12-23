package audio

import (
	"github.com/adrmcintyre/poweraid/audio/audiofilter"
	ebiten_audio "github.com/hajimehoshi/ebiten/v2/audio"
)

// Some alternate sound effects will be used when this is set
const alternateMode = false

// We simulate 3-channel audio
const channelCount = 3

// A Command represents a pending command to update
// the frequency and volume of an audio channel.
type Command struct {
	freq uint32
	vol  byte
}

// Audio encapsulates all of the audio state, and implements
// the io.Reader interface necessary for ebiten to be able to
// stream from it.
type Audio struct {
	// host playback
	player       *ebiten_audio.Player
	outputVolume int  // host output volume
	muted        bool // is all audio muted?

	// simulated hardware state
	hwVoice [voiceCount]hwVoice // current state of the simulated hardware
	pos     int64               // number of samples emitted into the stream
	filter  audiofilter.Filter  // output filter

	// audio sequencing
	nextSequence    int64                          // when we should next schedule the audio sequencer
	songProcessor   [channelCount]*SongProcessor   // state of the song processors
	effectProcessor [channelCount]*EffectProcessor // state of the effects processors
	soundCounter    byte                           // 60Hz counter for envelope effects, we're happy for it to wrap at 255
	command         [channelCount]Command          // pending output changes (referenced by the processors)
}

// NewAudio constructs an initialised Audio struct.
func NewAudio() *Audio {
	au := &Audio{}
	for i := range channelCount {
		channel := &au.command[i]
		au.songProcessor[i] = &SongProcessor{
			index:   i,
			counter: au,
			command: channel,
		}
		au.effectProcessor[i] = &EffectProcessor{
			index:   i,
			counter: au,
			command: channel,
		}
	}
	return au
}

// GetCounter implements the 'counter' interface for the benefit of the envelope generator.
func (au *Audio) GetCount() byte {
	return au.soundCounter
}

// PlaySong begins playing the specified song on channels 0 and 1.
// This takes precedence over any effects currently playing on
// those channels.
func (au *Audio) PlaySong(song SongId) {
	au.songProcessor[0].queueMask |= 1 << song // melody
	au.songProcessor[1].queueMask |= 1 << song // rhythm
}

// PlayTransientEffect plays a non-looping effect on channel 0.
// Any existing effects will continue to play.
func (au *Audio) PlayTransientEffect(i TransientEffectId) {
	au.effectProcessor[0].queueMask |= (1 << i)
}

// StopAllTransientEffects stops all active effects on channel 0.
func (au *Audio) StopAllTransientEffects() {
	au.effectProcessor[0].queueMask = 0
}

// PlayBackgroundEffect plays a looping effect on channel 1.
// The new effect replaces any background effects currently
// playing, with the exception of EnergiserEaten and EyesReturning
// which will continue to play.
func (au *Audio) PlayBackgroundEffect(i BackgroundEffectId) {
	const backgroundMask = (1 << EnergiserEaten) | (1 << EyesReturning)
	e := au.effectProcessor[1]
	background := e.queueMask & backgroundMask
	e.queueMask = (e.queueMask & background) | (1 << i)
}

// StopBackgroundEffect stops only the specified background effect.
func (au *Audio) StopBackgroundEffect(i BackgroundEffectId) {
	au.effectProcessor[1].queueMask &= ^(1 << i)
}

// StopAllBackgroundEffects stops all active background effects on channel 1.
func (au *Audio) StopAllBackgroundEffects() {
	au.effectProcessor[1].queueMask = 0
}

// PlayPacmanEffect plays a non-looping effect on channel 2.
// These effects all relate to pacman's activities.
// Any existing effects continue to play - except DotEatenEven
// or DotEvenOdd which will silence their opposite number.
func (au *Audio) PlayPacmanEffect(i PacmanEffectId) {
	const even = DotEatenEven
	const odd = DotEatenOdd
	const evenOddMask = byte(1<<even | 1<<odd)
	e := au.effectProcessor[2]
	if i == even || i == odd {
		qm := effectChannel[2].queueMask
		e.queueMask = (qm & ^evenOddMask) | (1 << i)
	} else {
		e.queueMask = (1 << i)
	}
}

// StopAllPacmanEffects stops all active effects on channel 2.
func (au *Audio) StopAllPacmanEffects() {
	au.effectProcessor[2].queueMask = 0
}
