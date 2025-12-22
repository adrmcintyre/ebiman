package data

// A PCM is notionally an infinitely repeating series of bits representing a
// pulse train.
//
// At each update, an actor advances every time a bit is shifted out of the top
// of the pulse train. Thus, the more bits set, the larger the proportion of
// updates when it advances, and the faster it moves.
type PCM uint32

// The bits in each pulse train are spread out as evenly as possible to reduce
// jerkiness. Each constant is named as a percentage of pacman's maximum speed.
const (
	PCM_5   PCM = 0x20000000 // 1 pulse every 32 updates
	PCM_10  PCM = 0x20002000 // 2
	PCM_15  PCM = 0x20040100 // 3
	PCM_20  PCM = 0x20202020 // 4
	PCM_25  PCM = 0x20810408 // 5
	PCM_30  PCM = 0x20842084 // 6
	PCM_35  PCM = 0x22110884 // 7
	PCM_40  PCM = 0x22222222 // 8
	PCM_45  PCM = 0x24489122 // 9
	PCM_50  PCM = 0x24922492 // 10
	PCM_55  PCM = 0x24924925 // 11
	PCM_60  PCM = 0x25252525 // 12
	PCM_65  PCM = 0x25A4925A // 13
	PCM_70  PCM = 0x259A259A // 14
	PCM_75  PCM = 0x2AAA5555 // 15
	PCM_80  PCM = 0x55555555 // 16
	PCM_85  PCM = 0x6AAAD555 // 17
	PCM_90  PCM = 0x6AD56AD5 // 18
	PCM_95  PCM = 0x5AD6B5AD // 19
	PCM_100 PCM = 0x6D6D6D6D // 20 - pacman's maximum speed
	PCM_105 PCM = 0x6DB6DB6D // 21
	PCM_110 PCM = 0x6DBB6DBB // 22
	PCM_MAX PCM = 0xFFFFFFFF // eyes return at full pelt
)

// Pulse rotates the pulse train by one bit, and returns
// true if the top bit was set.
func (pcm *PCM) Pulse() bool {
	msb := *pcm >> 31
	*pcm = (*pcm << 1) | msb
	return msb != 0
}
