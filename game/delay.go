package game

import "github.com/adrmcintyre/ebiman/data"

// ScheduleDelay starts a timer to delay the game state
// advancing for delayMilli milliseconds, causing all
// on screen action to pause.
func (g *Game) ScheduleDelay(delayMillis int) {
	g.DelayTimer = delayMillis * data.FPS / 1000
}

// CheckDelay checks if a delay is currently pending,
// in which case it renders the next frame without advancing
// the game state, and returns true.
func (g *Game) CheckDelay() bool {
	if g.DelayTimer <= 0 {
		return false
	}

	g.DelayTimer -= 1
	g.RenderFrameUncounted()
	return true
}
