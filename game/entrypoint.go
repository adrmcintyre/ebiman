package game

import (
	"time"

	"github.com/adrmcintyre/poweraid/audio"
	"github.com/adrmcintyre/poweraid/bonus"
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/ghost"
	"github.com/adrmcintyre/poweraid/input"
	"github.com/adrmcintyre/poweraid/level"
	"github.com/adrmcintyre/poweraid/message"
	"github.com/adrmcintyre/poweraid/option"
	"github.com/adrmcintyre/poweraid/pacman"
	"github.com/adrmcintyre/poweraid/player"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
	"github.com/hajimehoshi/ebiten/v2"
	ebiten_audio "github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	ActionRunLoop
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
	TASK_GHOST_RETURN TaskId = iota
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
	ScreenWidth    int
	ScreenHeight   int
	StatusMsg      message.Id
	PlayerMsg      message.Id
	RunningGame    bool
	StartMenuIndex int
	Action         Action
	Anim           AnimState
	Delay          int
	Tasks          [MAX_TASKS]Task
	TaskCount      int
	Options        option.Options
	PlayerNumber   int // current player
	SavedPlayer    [2]player.SavedState
	LevelState     level.State
	LevelConfig    level.Config
	Pacman         *pacman.Actor
	Ghosts         [4]*ghost.Actor
	BonusActor     bonus.Actor
	Video          video.Video
	Audio          *audio.Audio
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
		case TASK_GHOST_RETURN:
			// TODO - the only task - incorporate the delay processing here too?
			g.GhostReturn(task.Param)
		}
		copy(g.Tasks[:MAX_TASKS-1], g.Tasks[1:])
		g.TaskCount -= 1
	}

	switch g.Action {
	case ActionReset:
		g.ResetGame()
		g.Action = ActionSplash

	case ActionSplash:
		g.Anim.Frame = 0
		g.Anim.Delay = 0
		g.Action = ActionSplashAnim

	case ActionSplashAnim:
		if input.GetJoystickSwitch() {
			g.Anim.Animator = nil
			g.Action = ActionRunLoop
			g.RunningGame = false
		} else if g.Anim.Delay > 0 {
			g.Anim.Delay -= 1
		} else if frame, delay := g.SplashScreen(g.Anim.Frame); frame > 0 {
			g.Anim.Delay = delay
			g.Anim.Frame = frame
		} else {
			g.Action = ActionRunLoop
		}

	case ActionRunLoop:
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenWidth, g.ScreenHeight
}

func EntryPoint(w, h int) error {
	// "power on" sequence
	tile.MakeImages()
	sprite.MakeImages()
	color.MakeColorMatrixes()

	audioStream := audio.NewAudio()
	audioContext := ebiten_audio.NewContext(audio.SampleRate)
	audioPlayer, err := audioContext.NewPlayer(audioStream)
	if err != nil {
		return err
	}
	defer audioPlayer.Close()
	audioPlayer.SetBufferSize(100 * time.Millisecond)
	audioPlayer.Play()

	g := NewGame(w, h, audioStream)
	return ebiten.RunGame(g)
}

func NewGame(w, h int, audioStream *audio.Audio) *Game {
	pacman := pacman.NewActor()
	// ghosts are aware of pacman, and inky is also aware of blinky
	blinky := ghost.NewBlinky(pacman)
	pinky := ghost.NewPinky(pacman)
	inky := ghost.NewInky(pacman, blinky)
	clyde := ghost.NewClyde(pacman)

	return &Game{
		ScreenWidth:  w,
		ScreenHeight: h,
		Action:       ActionReset,
		Options:      option.MakeOptions(),
		PlayerNumber: 0,
		LevelState:   level.DefaultState(),
		LevelConfig:  level.DefaultConfig(),
		Pacman:       pacman,
		Ghosts:       [4]*ghost.Actor{blinky, pinky, inky, clyde},
		BonusActor:   bonus.MakeActor(),
		Video:        video.Video{},
		Audio:        audioStream,
	}
}
