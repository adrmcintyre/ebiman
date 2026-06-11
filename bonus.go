package main

import (
	"math/rand"

	"github.com/adrmcintyre/ebiman/audio"
	"github.com/adrmcintyre/ebiman/data"
)

// DropBonus begins the display of a bonus item.
// A timeout is established after which the bonus will be removed.
func (g *Game) DropBonus() {
	g.BonusActor.Visible = true
	// TODO should this be updates instead?
	minTime := data.MinBonusTime
	rangeTime := data.MaxBonusTime - minTime
	timeout := minTime + rand.Intn(rangeTime)
	g.Level.BonusTimeout = g.Level.FrameCounter + timeout
}

// HideBonus removes any bonus currently on display.
func (g *Game) HideBonus() {
	g.Level.BonusTimeout = 0
	g.BonusActor.Visible = false
}

// CheckTimeoutBonus checks if an active bonus has timed out,
// hiding it if necessary.
func (g *Game) CheckTimeoutBonus() {
	timeout := g.Level.BonusTimeout
	if timeout != 0 && g.Level.FrameCounter >= timeout {
		g.HideBonus()
	}
}

// EatBonus consumes the current bonus.
func (g *Game) EatBonus() {
	g.IncrementScore(g.LevelConfig.BonusInfo.Score)

	g.Player.BonusStatus.AddBonus(g.LevelConfig.BonusType)
	g.Level.BonusTimeout = 0
	g.Level.BonusScoreTimeout = g.Level.FrameCounter + 2*data.FPS

	g.Audio.PlayPacmanEffect(audio.FruitEaten)

	g.BonusActor.Visible = false
}

// HideBonusScore hides any currently visible bonus score text.
func (g *Game) HideBonusScore() {
	g.Level.BonusScoreTimeout = 0
}

// CheckTimeoutBonusScore checks is an active bonus-score has timed out,
// hiding it if necessary.
func (g *Game) CheckTimeoutBonusScore() {
	timeout := g.Level.BonusScoreTimeout
	if timeout > 0 && g.Level.FrameCounter >= timeout {
		g.HideBonusScore()
	}
}
