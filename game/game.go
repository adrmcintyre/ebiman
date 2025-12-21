package game

import (
	"time"

	"github.com/adrmcintyre/poweraid/audio"
	"github.com/adrmcintyre/poweraid/bonus"
	"github.com/adrmcintyre/poweraid/color"
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
)

const (
	quitKey = ebiten.KeyQ
)

// A Game collects all state related to the running of the game.
type Game struct {
	ScreenWidth  int          // width of screen in logical pixels
	ScreenHeight int          // height of screen in logical pixels
	Video        *video.Video // simulated video hardware
	Audio        *audio.Audio // simulated audio hardware

	DelayTimer     int       // delay timer in frames (if non-zero)
	TaskQueue      []Task    // pending tasks to execute
	GameState      GameState // current game state
	CoroState      CoroState // state of currently executing coroutine
	StartMenuIndex int       // currently selected menu item in options screen

	// core game state
	RunningGame  bool                 // is the game core loop in progress?
	Options      option.Options       // game options
	PlayerNumber int                  // current player, 0 or 1
	SavedPlayer  [2]player.SavedState // saved states of each player
	LevelState   level.State          // state of level in progress
	LevelConfig  level.Config         // configuration of current level

	// in-game prompts
	StatusMsg message.Id // possible status message in maze (ready / game over)
	PlayerMsg message.Id // ppossible layer message in maze (player 1 / 2)

	// the actors
	Pacman     *pacman.Actor   // pacman's state
	Ghosts     [4]*ghost.Actor // each ghost's state
	BonusActor *bonus.Actor    // the bonus's state
}

// NewGame returns a default-initialised Game object.
func NewGame(w, h int) *Game {
	pacman := pacman.NewActor()

	// ghosts are aware of pacman, and inky is also aware of blinky
	blinky := ghost.NewBlinky(pacman)
	pinky := ghost.NewPinky(pacman)
	inky := ghost.NewInky(pacman, blinky)
	clyde := ghost.NewClyde(pacman)

	bonusActor := bonus.NewActor()

	return &Game{
		ScreenWidth:  w,
		ScreenHeight: h,
		GameState:    GameStateReset,
		Options:      option.DefaultOptions(),
		PlayerNumber: 0,
		LevelState:   level.DefaultState(),
		LevelConfig:  level.DefaultConfig(),
		Pacman:       pacman,
		Ghosts:       [4]*ghost.Actor{blinky, pinky, inky, clyde},
		BonusActor:   bonusActor,
		// hook these up later
		Video: nil,
		Audio: nil,
	}
}

// Execute initialises the "hardware" and begins the execution
// of the game via the ebiten framework.
func (g *Game) Execute() error {
	// pre-compute static data
	tile.MakeImages()
	sprite.MakeImages()
	color.MakeColorMatrixes()

	// hookup video "hardware"
	g.Video = &video.Video{}

	// hookup audio "hardware"
	g.Audio = audio.NewAudio()
	audioContext := ebiten_audio.NewContext(audio.SampleRate)
	audioPlayer, err := audioContext.NewPlayer(g.Audio)
	if err != nil {
		return err
	}
	defer audioPlayer.Close()

	// If the audio buffer size is too long, audio will lag the action;
	// if it's too short, the audio becomes choppy.
	audioPlayer.SetBufferSize(100 * time.Millisecond)
	audioPlayer.Play()

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
	if input.QuitPressed() {
		return ebiten.Termination
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
}

// Layout is called by the ebiten framework to establish the size of the
// rendered image.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenWidth, g.ScreenHeight
}

// assert that Game implements the ebiten.Game interface.
var _ ebiten.Game = (*Game)(nil)
