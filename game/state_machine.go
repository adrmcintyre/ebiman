package game

import (
	"github.com/adrmcintyre/ebiman/input"
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
		// splash screen coro gets special handling
		//
		// TODO it would be really nice to unify this code
		// with the normal coroutine handler below.
		// Probably by creating a method to set g.GameState,
		// and setting it as the next continuation.
		g.RenderFrame()
		if input.GetJoystickSwitch() {
			g.Coro = nil
			g.GameState = GameStateCoreLoop
			g.RunningGame = false
		} else if !g.Coro.invoke(g) {
			g.Coro = nil
			g.GameState = GameStateCoreLoop
		}

	case GameStateCoreLoop:
		g.RenderFrame()
		{
		updateTwice:
			for range 2 {
				var ret Return
				if g.Coro != nil {
					next := g.Coro.next
					if !g.Coro.invoke(g) {
						g.Coro = nil
						ret = next(g)
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
