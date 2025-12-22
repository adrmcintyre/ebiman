package audio

import (
	"io"
	"time"

	ebiten_audio "github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	// Specifies the audio system's sample frequency
	sampleRate = 48000

	// Configures the audio buffer size.
	// If it's too long, audio will lag the action.
	// If it's too short, the audio becomes choppy.
	bufferSize = 60 * time.Millisecond
)

// Configure some parameters of the simulated hardware
const (
	voiceCount  = 3               // how many voices are supported
	volumeCount = 16              // how many distinct volume levels in a sample
	maxVolume   = volumeCount - 1 // the maximum sample volume
	encodedZero = uint16(0x8000)  // how is zero-output encoded in the output stream
)

// An hwVoice represents the current state of the registers for 1 voice.
type hwVoice struct {
	wave byte   // low 3 bits used – selects waveform 0-7 from ROM
	vol  byte   // low nibble – 0 off to 15 loudest
	freq uint32 // real hardware has 20 bits for voice 0; 16 bits voices 1, 2
}

// NewPlayer returns an ebiten audio player configured
// to play the output of the simulate hardware.
func (au *Audio) NewPlayer() (*ebiten_audio.Player, error) {
	audioContext := ebiten_audio.NewContext(sampleRate)
	audioPlayer, err := audioContext.NewPlayer(au)
	if err == nil {
		audioPlayer.SetBufferSize(bufferSize)
		audioPlayer.Play()
	}
	return audioPlayer, err
}

// Read is io.Reader's Read.
//
// Read fills buf with sampled audio according to hwVoice settings.
func (au *Audio) Read(buf []byte) (int, error) {
	const (
		bytesPerValue  = 2
		bytesPerSample = bytesPerValue * 2 // 2 x 16-bit samples (for left and right)
	)

	alignedLen := len(buf) / bytesPerSample * bytesPerSample

	numSamples := alignedLen / bytesPerSample
	numEmitted := au.pos / bytesPerSample

	// sample and mix channels
	for i := range numSamples {
		sampleIndex := float64(numEmitted + int64(i))
		t := sampleIndex / sampleRate

		// schedule the sequencer every 1/60s
		if t >= au.nextFrameTime {
			au.nextFrameTime = t + 1.0/60.0
			au.nudgeSequencer()
		}

		v16 := encodedZero
		if !au.muted {
			for ch, channel := range au.hwVoice {
				freq := channel.freq * 3
				// channel 0 has more freq bits allocated
				if ch == 0 {
					freq >>= 4
				}
				j := int(lookupCount*float64(freq)/2*t) % lookupCount
				v16 += scaledWaveData[channel.vol][channel.wave][j]
			}
		}

		// encode left channel
		buf[4*i] = byte(v16)
		buf[4*i+1] = byte(v16 >> 8)
		// same audio on the right channel
		buf[4*i+2] = byte(v16)
		buf[4*i+3] = byte(v16 >> 8)
	}

	au.pos += int64(alignedLen)

	return alignedLen, nil
}

// Close is io.Closer's Close.
func (au *Audio) Close() error {
	return nil
}

// assert we implement the interface
var _ io.Closer = (*Audio)(nil)
