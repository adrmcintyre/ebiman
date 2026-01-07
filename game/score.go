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
	if oldScore < data.ExtraLifeScore && newScore >= data.ExtraLifeScore {
		ls.AwardExtraLife()
		g.Audio.PlayTransientEffect(audio.ExtraLife)

	}

	ls.SetScore(playerNumber, newScore)
}

func (g *Game) RegisterScore() {
	score := g.LevelState.Score1
	if g.PlayerNumber == 1 {
		score = g.LevelState.Score2
	}

	lb := g.Options.LeaderboardName()
	g.Service.RegisterScore(lb, int64(score))
}
