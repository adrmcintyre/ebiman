package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/input"
	"github.com/adrmcintyre/poweraid/message"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

type ReturnAction int

const (
	rStop ReturnAction = iota
	rContinue
	rWithAnim
)

type Animator func(*Game, int) (int, int)
type Continuation func(*Game) Return

type Return struct {
	ra   ReturnAction
	anim Animator
	next Continuation
}

var ThenStop = Return{rStop, nil, nil}
var ThenContinue = Return{rContinue, nil, nil}

func WithAnim(anim Animator, next Continuation) Return {
	return Return{rWithAnim, anim, next}
}

type Action int

const (
	ActionReset Action = iota
	ActionSplash
	ActionSplashAnim
	ActionReady
	ActionRun
)

type AnimState struct {
	Animator Animator
	Frame    int
	Delay    int
	Next     Continuation
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
	StatusMsg      message.MsgId
	PlayerMsg      message.MsgId
	RunningGame    bool
	StartMenuIndex int
	Action         Action
	Anim           AnimState
	Delay          int
	Tasks          [MAX_TASKS]Task
	TaskCount      int
	Options        GameOptions
	PlayerNumber   int // current player
	SavedPlayer    [2]SavedPlayerState
	LevelState     LevelState
	LevelConfig    LevelConfig
	Pacman         PacmanActor
	Ghosts         [4]GhostActor
	BonusActor     BonusActor
	Video          video.Video
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if g.Delay > 0 {
		g.Delay -= 1
		g.RenderFrameUncounted()
		return nil
	}

	for g.TaskCount > 0 {
		switch task := g.Tasks[0]; task.Id {
		case TaskReturnGhost:
			// TODO - the only task - incorporate the delay processing here too?
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
		if input.GetJoystickSwitch() {
			g.Action = ActionRun
		} else if g.Anim.Delay > 0 {
			g.Anim.Delay -= 1
		} else if frame, delay := g.SplashScreen(g.Anim.Frame); frame > 0 {
			g.Anim.Delay = delay
			g.Anim.Frame = frame
		} else {
			g.Action = ActionRun
		}

	case ActionRun:
		g.RenderFrame()
		{
		updateLoop:
			for range 2 {
				var ret Return
				if animator := g.Anim.Animator; animator != nil {
					if g.Anim.Delay > 0 {
						g.Anim.Delay -= 1
					} else if frame, delay := animator(g, g.Anim.Frame); frame > 0 {
						g.Anim.Delay = delay
						g.Anim.Frame = frame
					} else {
						g.Anim.Animator = nil
						ret = g.Anim.Next(g)
					}
				} else {
					ret = g.UpdateState()
				}
				switch ret.ra {
				case rStop:
					break updateLoop
				case rContinue:
				case rWithAnim:
					g.Anim.Animator = ret.anim
					g.Anim.Frame = 0
					g.Anim.Delay = 0
					g.Anim.Next = ret.next
				}
			}
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
		Video:        video.Video{},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
