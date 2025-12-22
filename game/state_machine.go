package game

import (
	"github.com/adrmcintyre/poweraid/input"
)

// A GameState identifies a state in the games's top-level state machine.
type GameState int

const (
	GameStateReset        GameState = iota // game is resetting (power on)
	GameStateSplashStart                   // splash screen is to be displayed
	GameStateSplashScreen                  // splash screen is running
	GameStateCoreLoop                      // core game loop is running
)

// RunStateMachine executes the current game action.
func (g *Game) RunStateMachine() {
	switch g.GameState {
	case GameStateReset:
		g.ResetGame()
		g.GameState = GameStateSplashStart

	case GameStateSplashStart:
		g.Coro = &Coro{method: (*Game).SplashScreen}
		g.GameState = GameStateSplashScreen

	case GameStateSplashScreen:
		if input.GetJoystickSwitch() {
			g.Coro = nil
			g.GameState = GameStateCoreLoop
			g.RunningGame = false
		} else {
			g.Coro = g.Coro.invoke(g)
			if g.Coro == nil {
				g.GameState = GameStateCoreLoop
			}
		}

	case GameStateCoreLoop:
		g.RenderFrame()
		{
		updateTwice:
			for range 2 {
				var ret Return
				if coro := g.Coro; coro != nil {
					if coro.invoke(g) == nil {
						g.Coro = nil
						ret = coro.next(g)
					}
				} else {
					ret = g.UpdateState()
				}

				if ret.coro != nil {
					g.Coro = ret.coro
				}
				if ret.done {
					break updateTwice
				}
			}
		}
	}
}
