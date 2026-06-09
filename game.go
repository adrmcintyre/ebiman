package main

import (
	"github.com/adrmcintyre/ebiman/actor"
	"github.com/adrmcintyre/ebiman/audio"
	"github.com/adrmcintyre/ebiman/input"
	"github.com/adrmcintyre/ebiman/platform"
	"github.com/adrmcintyre/ebiman/service"
	"github.com/adrmcintyre/ebiman/state"
	"github.com/adrmcintyre/ebiman/video"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	hTiles, vTiles        = 28, 36 // dimensions of display area in tiles
	tileWidth, tileHeight = 8, 8   // dimensions of a tile in simulated pixels
	border                = 8      // small border around display in simulated pixels

	// calculate logical output size
	logicalWidth  = float64(hTiles*tileWidth + 2*border)
	logicalHeight = float64(vTiles*tileHeight + 2*border)
	logicalAspect = float64(logicalWidth) / float64(logicalHeight)
)

// A Game collects all state related to the running of the game.
type Game struct {
	Video       *video.Video     // simulated video hardware
	Audio       *audio.Audio     // simulated audio hardware
	Input       *input.Input     // user input
	Service     *service.Service // game service client
	IsWasmBuild bool             // are we running in wasm?

	DelayTimer     int       // delay timer in frames (if non-zero)
	TaskQueue      []Task    // pending tasks to execute
	GameState      GameState // current game state
	Coro           *Coro     // currently executing coroutine, if non-nil
	StartMenuIndex int       // currently selected menu item in options screen
	Paused         bool      // currently paused

	// core game state
	RunningGame  bool                      // is the game core loop in progress?
	Options      Options                   // game options
	PlayerNumber int                       // current player, 0 or 1
	SavedPlayer  [2]state.SavedPlayerState // saved states of each player
	LevelState   LevelState                // state of level in progress
	LevelConfig  LevelConfig               // configuration of current level

	// in-game prompts
	StatusMsg MsgId // possible status message in maze (ready / game over)
	PlayerMsg MsgId // possible layer message in maze (player 1 / 2)

	// the actors
	Pacman     *actor.Pacman   // pacman's state
	Ghosts     [4]*actor.Ghost // each ghost's state
	BonusActor *actor.Bonus    // the bonus's state
}

// NewGame returns a default-initialised Game object.
func NewGame(serverUrl string, serverKey string, isWasmBuild bool) *Game {
	pacman := actor.NewPacman()

	// ghosts are aware of pacman, and inky is also aware of blinky
	blinky := actor.NewBlinky(pacman)
	pinky := actor.NewPinky(pacman)
	inky := actor.NewInky(pacman, blinky)
	clyde := actor.NewClyde(pacman)

	bonusActor := actor.NewBonus()

	inp := input.New()
	if isWasmBuild {
		inp.SetTouchLayout(input.MakeTouchLayout(input.LayoutRectsLRUD, 0, int(logicalHeight), int(logicalWidth), 120))
	}

	return &Game{
		Video:       nil,
		Audio:       nil,
		Input:       inp,
		Service:     service.New(serverUrl, serverKey),
		IsWasmBuild: isWasmBuild,

		GameState: GameStateReset,

		Options:      DefaultOptions(),
		PlayerNumber: 0,
		LevelState:   DefaultLevelState(),
		LevelConfig:  DefaultLevelConfig(),

		Pacman:     pacman,
		Ghosts:     [4]*actor.Ghost{blinky, pinky, inky, clyde},
		BonusActor: bonusActor,
	}
}

// Execute initialises the "hardware" and begins the execution
// of the game via the ebiten framework.
func (g *Game) Execute() error {
	// pre-compute static data
	video.InitTiles()
	video.InitSprites()
	video.InitColors()

	// hookup video "hardware"
	g.Video = &video.Video{}
	if err := g.Video.Init(); err != nil {
		return err
	}

	var latency audio.Latency
	if g.IsWasmBuild {
		latency = audio.LatencyHigh
	} else {
		latency = audio.LatencyLow
	}

	// hookup audio "hardware"
	g.Audio = audio.NewAudio()
	defer func() { _ = g.Audio.Close() }()

	// connect to host's audio
	if err := g.Audio.Connect(latency); err != nil {
		return err
	}

	if deviceId, ok := platform.GetDeviceId(); ok {
		g.Service.Auth(deviceId)
	}

	return ebiten.RunGame(g)
}

// Update is called by the ebiten framework to progress the game state.
//
// Here we both progress the game state, and "render" the next frame.
// However we are only rendering into the simulated hardware's video
// memory by writing tiles and sprite descriptors.
//
// The real user-facing rendering only happens when Draw() is called
// by the framework.
func (g *Game) Update() error {
	g.Input.Update()

	if !g.IsWasmBuild && g.Input.Quit() {
		return ebiten.Termination
	}
	if g.Input.VolumeUp() {
		g.Audio.OutputVolumeUp()
	}
	if g.Input.VolumeDown() {
		g.Audio.OutputVolumeDown()
	}
	if g.Input.Pause() {
		g.Paused = !g.Paused
	}

	if g.Paused {
		return nil
	}
	if g.CheckDelay() {
		return nil
	}
	g.RunTaskQueue()
	g.RunStateMachine()

	return nil
}

// Draw is called by the ebiten framework to prepare the next frame.
//
// We paint the simulated hardware's video buffer into the supplied bitmap.
func (g *Game) Draw(screen *ebiten.Image) {
	g.Video.Draw(screen)
	// render any input controls
	g.Input.Draw(screen)
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("fps=%.1f tps=%.1f", ebiten.ActualFPS(), ebiten.ActualTPS()), 0, 288)
}

// Layout is called by the ebiten framework to establish the size of the
// rendered image.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	var (
		fOutsideWidth  = float64(outsideWidth)
		fOutsideHeight = float64(outsideHeight)
		outputAspect   = fOutsideWidth / fOutsideHeight

		fScreenWidth  = logicalWidth
		fScreenHeight = logicalHeight
		scale         float64
	)

	// centre output in window
	if outputAspect > logicalAspect {
		scale = fOutsideHeight / logicalHeight
		fScreenWidth = logicalHeight * outputAspect
	} else {
		scale = fOutsideWidth / logicalWidth
		fScreenHeight = logicalWidth / outputAspect
	}
	g.Video.SetOffset(
		outsideWidth-int(logicalWidth*scale),
		0, //outsideHeight-int(logicalHeight*scale),
	)
	return int(fScreenWidth), int(fScreenHeight)
}

// assert that Game implements the ebiten.Game interface.
var _ ebiten.Game = (*Game)(nil)

// DrawFinalScreen is called by the ebiten framework to apply effects to
// the final rendered output.
func (g *Game) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	g.Video.PostProcess(screen, offscreen)
}

var _ ebiten.FinalScreenDrawer = (*Game)(nil)
