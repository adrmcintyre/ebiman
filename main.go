package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/tile"
)

type Action int

const (
	ActionReset Action = iota
	ActionSplash
	ActionSplashAnim
	ActionMain
	ActionReady
	ActionRun
)

type AnimState struct {
	Frame int
	Delay int
}

const MAX_TASKS = 16

type TaskId int

const (
	TaskReturnGhost TaskId = iota
)

type Task struct {
	Id    TaskId
	Param int
}

func (g *Game) AddTask(id TaskId, param int) {
	g.Tasks[g.TaskCount] = Task{id, param}
	g.TaskCount += 1
}

func (g *Game) ScheduleDelay(delay int) {
	g.Delay = delay * data.FPS / 1000
}

type Game struct {
	Action       Action
	Anim         AnimState
	Delay        int
	Tasks        [MAX_TASKS]Task
	TaskCount    int
	Options      GameOptions
	PlayerNumber int // current player
	SavedPlayer  [2]SavedPlayerState
	LevelState   LevelState
	LevelConfig  LevelConfig
	Pacman       PacmanActor
	Ghosts       [4]GhostActor
	BonusActor   BonusActor
	Video        Video
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if g.Delay > 0 {
		g.Delay -= 1
		g.RenderFrame()
		//g.LevelState.FrameCounter += 1
		return nil
	}

	for g.TaskCount > 0 {
		switch task := g.Tasks[0]; task.Id {
		case TaskReturnGhost:
			g.ReturnGhost(task.Param)
		}
		copy(g.Tasks[:MAX_TASKS-1], g.Tasks[1:])
		g.TaskCount -= 1
	}

	switch g.Action {
	case ActionReset:
		g.ResetGame()

	case ActionSplash:
		g.Anim.Frame = 0
		g.Anim.Delay = 0
		g.Action = ActionSplashAnim

	case ActionSplashAnim:
		if g.Anim.Delay > 0 {
			g.Anim.Delay -= 1
		} else if frame, delay := g.SplashScreen(g.Anim.Frame); frame > 0 {
			g.Anim.Delay = delay
			g.Anim.Frame = frame
		} else {
			g.Action = ActionMain
		}

	case ActionMain:
		g.MainGame()

	case ActionRun:
		g.RenderFrame()
		g.LevelState.FrameCounter += 1
		for range 2 {
			if !g.UpdateState() {
				break
			}
			g.LevelState.UpdateCounter += 1
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Video.Draw(screen)
}

const (
	hBorder     = 8
	vBorder     = 8
	hWidth      = 224
	vHeight     = 288
	totalWidth  = hWidth + 2*hBorder
	totalHeight = vHeight + 2*vBorder
	screenScale = 2.5
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return totalWidth, totalHeight
}

func main() {
	ebiten.SetWindowSize(int(totalWidth*screenScale), int(totalHeight*screenScale))
	ebiten.SetWindowTitle("PowerAid")

	tile.MakeImages()
	sprite.MakeImages()
	palette.MakeColorMatrixes()

	game := &Game{
		Action:       ActionReset,
		Options:      DefaultGameOptions(),
		PlayerNumber: 0,
		LevelState:   DefaultLevelState(),
		LevelConfig:  DefaultLevelConfig(),
		Pacman:       MakePacman(),
		Ghosts:       MakeGhosts(),
		BonusActor:   MakeBonus(),
		Video:        Video{},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
