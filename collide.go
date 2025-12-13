package main

import (
	"math/rand"

	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/tile"
)

func (g *Game) EatGhost(ghost *GhostActor) {
	ghostScore := &data.GhostScore[g.LevelState.GhostsEaten]
	g.LevelState.IncrementScore(&g.Video, g.PlayerNumber, ghostScore.Score)

	ghost.ScoreSprite = ghostScore.Sprite

	g.Pacman.Motion.Visible = false

	g.RenderFrame()

	ghost.Mode = MODE_RETURNING
	ghost.Motion.Pcm = data.PCM_MAX

	g.ScheduleDelay(data.DISPLAY_GHOST_SCORE_MS)
	g.AddTask(TaskReturnGhost, ghost.Id)
}

func (g *Game) ReturnGhost(id int) {
	ghost := &g.Ghosts[id]
	ghost.ScoreSprite = 0

	g.Pacman.Motion.Visible = true

	g.LevelState.GhostsEaten += 1
}

func (g *Game) EatPill() {
	g.LevelState.IncrementScore(&g.Video, g.PlayerNumber, data.DOT_SCORE)
	g.CountDot()
	g.Pacman.StallTimer = data.DOT_STALL
}

func (g *Game) DropBonus() {
	g.BonusActor.Motion.Visible = true
	// TODO should this be updates instead?
	minTime := data.MIN_BONUS_TIME
	rangeTime := data.MAX_BONUS_TIME - minTime
	timeout := minTime + rand.Intn(rangeTime)
	g.LevelState.BonusTimeout = g.LevelState.FrameCounter + timeout
}

func (g *Game) HideBonus() {
	g.LevelState.BonusTimeout = 0
	g.BonusActor.Motion.Visible = false
}

func (g *Game) TimeoutBonus() {
	timeout := g.LevelState.BonusTimeout
	if timeout != 0 && g.LevelState.FrameCounter >= timeout {
		g.HideBonus()
	}
}

func (g *Game) EatBonus() {
	g.LevelState.IncrementScore(&g.Video, g.PlayerNumber, g.LevelConfig.BonusInfo.Score)
	g.LevelState.BonusState.AddBonus(&g.Video, g.LevelConfig.BonusType)
	g.LevelState.BonusTimeout = 0
	g.LevelState.BonusScoreTimeout = g.LevelState.FrameCounter + 2*data.FPS

	g.BonusActor.Motion.Visible = false
	g.Video.SetCursor(12, 20)
	g.Video.WriteTiles(g.LevelConfig.BonusInfo.Tiles, palette.SCORE)
}

func (g *Game) HideBonusScore() {
	g.LevelState.BonusScoreTimeout = 0
	g.Video.SetCursor(12, 20)
	for range 4 {
		g.Video.WriteTile(tile.SPACE, palette.BLACK)
	}
}

func (g *Game) TimeoutBonusScore() {
	timeout := g.LevelState.BonusScoreTimeout
	if timeout > 0 && g.LevelState.FrameCounter >= timeout {
		g.HideBonusScore()
	}
}

func (g *Game) CountDot() {
	g.LevelState.DotsRemaining -= 1
	g.LevelState.DotsEaten += 1

	switch g.LevelState.DotsEaten {
	case data.FIRST_BONUS_DOTS, data.SECOND_BONUS_DOTS:
		g.DropBonus()
	}

	g.PacmanIsActive()

	if g.LevelState.GlobalDotCounterEnabled {
		g.LevelState.GlobalDotCounter += 1
	} else {
		for j := 1; j < 4; j++ {
			ghost := &g.Ghosts[j]
			if ghost.Mode == MODE_HOME {
				ghost.DotCounter += 1
				break
			}
		}
	}
}

func (g *Game) EatPower() {
	g.LevelState.IncrementScore(&g.Video, g.PlayerNumber, data.POWER_SCORE)
	g.CountDot()
	g.Pacman.StallTimer = data.POWER_STALL
	g.Pacman.Motion.Pcm = g.LevelConfig.Speeds.PacmanBlue

	g.LevelState.BlueTimeout = g.LevelState.UpdateCounter + g.LevelConfig.BlueTime
	g.LevelState.WhiteBlueTimeout = g.LevelState.BlueTimeout - g.LevelConfig.WhiteBlueCount*data.WHITE_BLUE_PERIOD
	g.LevelState.IsFlashing = false
	g.LevelState.IsWhite = false
	g.LevelState.GhostsEaten = 0

	// If some ghost is already scared, don't scare additional ghosts
	alreadyScared := false
	for j := range 4 {
		if g.Ghosts[j].SubMode == SUBMODE_SCARED {
			alreadyScared = true
			break
		}
	}

	if !alreadyScared {
		for j := range 4 {
			ghost := &g.Ghosts[j]
			if ghost.Mode == MODE_PLAYING || ghost.Mode == MODE_HOME {
				ghost.GhostSetSubmode(SUBMODE_SCARED)
				ghost.Motion.Pcm = g.LevelConfig.Speeds.GhostBlue
			}
		}
	}
}

// returns true if pacman dies
func (g *Game) CollidePacman() bool {
	v := &g.Video

	x := g.Pacman.Motion.Pos.X / 8
	y := g.Pacman.Motion.Pos.Y / 8

	t := v.GetTile(x, y)

	switch t {
	case tile.PILL:
		v.SetTile(x, y, tile.SPACE)
		g.EatPill()
	case tile.POWER, tile.POWER_SMALL:
		v.SetTile(x, y, tile.SPACE)
		g.EatPower()
	}

	if (g.LevelState.BonusTimeout > 0) &&
		(g.BonusActor.Motion.Pos.X/8 == x) &&
		(g.BonusActor.Motion.Pos.Y/8 == y) {
		g.EatBonus()
	}

	for j := range 4 {
		ghost := &g.Ghosts[j]
		if (ghost.Mode == MODE_PLAYING) &&
			(ghost.SubMode == SUBMODE_SCARED) &&
			(ghost.Motion.Pos.X/8 == x) &&
			(ghost.Motion.Pos.Y/8 == y) {
			g.EatGhost(ghost)
		}
	}

	for j := range 4 {
		ghost := &g.Ghosts[j]
		if (ghost.Mode == MODE_PLAYING) &&
			(ghost.SubMode != SUBMODE_SCARED) &&
			(ghost.Motion.Pos.X/8 == x) &&
			(ghost.Motion.Pos.Y/8 == y) {
			return true
		}
	}

	return false
}
