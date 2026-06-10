package main

import (
	"github.com/adrmcintyre/ebiman/audio"
	"github.com/adrmcintyre/ebiman/data"
)

// IncrementScore performs the necessary actions for
// increasing the current player's score, awarding an
// extra life when necessary.
func (g *Game) IncrementScore(delta int) {
	if g.DemoMode {
		return
	}

	oldScore := g.Player.Score
	newScore := oldScore + delta

	// pac man very generously awards one and only one extra life!
	if oldScore < data.ExtraLifeScore && newScore >= data.ExtraLifeScore {
		g.Player.AwardExtraLife()
		g.Audio.PlayTransientEffect(audio.ExtraLife)

	}

	g.SetScore(g.PlayerNumber, newScore)
}

// SetScore records the specified player's latest score,
// updating the highscore if appropriate.
func (g *Game) SetScore(playerNumber int, score int) {
	if score > g.HighScore {
		g.HighScore = score
	}
	g.Players[playerNumber].Score = score
}

// RegisterScore registers the current player's score with the leaderboard.
func (g *Game) RegisterScore() {
	score := g.Player.Score
	lb := g.Options.LeaderboardName()
	g.Service.RegisterScore(lb, int64(score))
}

// RefreshHighScore refreshes the high score from the leaderboard
// for the currently selected game mode.
func (g *Game) RefreshHighScore() {
	// TODO we need a cache in the game itself in case leaderboard is unavailable.
	lb := g.Options.LeaderboardName()
	highScore, ok := g.Service.GetHighScore(lb)
	if ok {
		g.HighScore = highScore
	}
}

// WriteScores writes the tiles for displaying the current
// high-score and player(s) scores into the top status area.
func (g *Game) WriteScores(numPlayers int) {
	v := g.Video
	v.WriteHighScore(g.HighScore)
	v.WriteScoreAt(1, 1, g.Players[0].Score)
	if numPlayers > 1 {
		v.WriteScoreAt(20, 1, g.Players[1].Score)
	}
}
