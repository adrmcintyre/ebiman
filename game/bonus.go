package game

import (
	"math/rand"

	"github.com/adrmcintyre/poweraid/data"
)

func (g *Game) DropBonus() {
	g.BonusActor.Visible = true
	// TODO should this be updates instead?
	minTime := data.MIN_BONUS_TIME
	rangeTime := data.MAX_BONUS_TIME - minTime
	timeout := minTime + rand.Intn(rangeTime)
	g.LevelState.BonusTimeout = g.LevelState.FrameCounter + timeout
}

func (g *Game) HideBonus() {
	g.LevelState.BonusTimeout = 0
	g.BonusActor.Visible = false
}

func (g *Game) TimeoutBonus() {
	timeout := g.LevelState.BonusTimeout
	if timeout != 0 && g.LevelState.FrameCounter >= timeout {
		g.HideBonus()
	}
}

func (g *Game) EatBonus() {
	g.LevelState.IncrementScore(g.PlayerNumber, g.LevelConfig.BonusInfo.Score)
	g.LevelState.BonusState.AddBonus(g.LevelConfig.BonusType)
	g.LevelState.BonusTimeout = 0
	g.LevelState.BonusScoreTimeout = g.LevelState.FrameCounter + 2*data.FPS

	g.BonusActor.Visible = false
}

func (g *Game) HideBonusScore() {
	g.LevelState.BonusScoreTimeout = 0
}

func (g *Game) TimeoutBonusScore() {
	timeout := g.LevelState.BonusScoreTimeout
	if timeout > 0 && g.LevelState.FrameCounter >= timeout {
		g.HideBonusScore()
	}
}
