package game

import (
	"math/rand"

	"github.com/adrmcintyre/ebiman/audio"
	"github.com/adrmcintyre/ebiman/bonus"
	"github.com/adrmcintyre/ebiman/data"
)

// DropBonus begins the display of a bonus item.
// A timeout is established after which the bonus will be removed.
func (g *Game) DropBonus() {
	g.BonusActor.Visible = true
	// TODO should this be updates instead?
	minTime := bonus.MinBonusTime
	rangeTime := bonus.MaxBonusTime - minTime
	timeout := minTime + rand.Intn(rangeTime)
	g.LevelState.BonusTimeout = g.LevelState.FrameCounter + timeout
}

// HideBonus removes any bonus currently on display.
func (g *Game) HideBonus() {
	g.LevelState.BonusTimeout = 0
	g.BonusActor.Visible = false
}

// CheckTimeoutBonus checks if an active bonus has timed out,
// hiding it if necessary.
func (g *Game) CheckTimeoutBonus() {
	timeout := g.LevelState.BonusTimeout
	if timeout != 0 && g.LevelState.FrameCounter >= timeout {
		g.HideBonus()
	}
}

// EatBonus consumes the current bonus.
func (g *Game) EatBonus() {
	g.IncrementScore(g.LevelConfig.BonusInfo.Score)

	g.LevelState.BonusStatus.AddBonus(g.LevelConfig.BonusType)
	g.LevelState.BonusTimeout = 0
	g.LevelState.BonusScoreTimeout = g.LevelState.FrameCounter + 2*data.FPS

	g.Audio.PlayPacmanEffect(audio.FruitEaten)

	g.BonusActor.Visible = false
}

// HideBonusScore hides any currently visible bonus score text.
func (g *Game) HideBonusScore() {
	g.LevelState.BonusScoreTimeout = 0
}

// CheckTimeoutBonusScore checks is an active bonus-score has timed out,
// hiding it if necessary.
func (g *Game) CheckTimeoutBonusScore() {
	timeout := g.LevelState.BonusScoreTimeout
	if timeout > 0 && g.LevelState.FrameCounter >= timeout {
		g.HideBonusScore()
	}
}
