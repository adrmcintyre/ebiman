package game

import (
	"github.com/adrmcintyre/ebiman/audio"
	"github.com/adrmcintyre/ebiman/data"
)

// IncrementScore performs the necessary actions for
// increasing the current player's score, awarding an
// extra live when necessary.
func (g *Game) IncrementScore(delta int) {
	ls := &g.LevelState
	if ls.DemoMode {
		return
	}

	playerNumber := g.PlayerNumber

	oldScore := ls.Score1
	if playerNumber == 1 {
		oldScore = ls.Score2
	}
	newScore := oldScore + delta

	// pac man very generously awards one and only one extra life!
	if oldScore < data.EXTRA_LIFE_SCORE && newScore >= data.EXTRA_LIFE_SCORE {
		ls.AwardExtraLife()
		g.Audio.PlayTransientEffect(audio.ExtraLife)

	}

	ls.SetScore(playerNumber, newScore)
}
