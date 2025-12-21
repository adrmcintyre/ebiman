package game

import "github.com/adrmcintyre/poweraid/input"

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
		g.CoroState.step = 0
		g.CoroState.delay = 0
		g.GameState = GameStateSplashScreen

	case GameStateSplashScreen:
		if input.GetJoystickSwitch() {
			g.CoroState.coro = nil
			g.GameState = GameStateCoreLoop
			g.RunningGame = false
		} else if g.CoroState.delay > 0 {
			g.CoroState.delay -= 1
		} else if frame, delay := g.SplashScreen(g.CoroState.step); frame > 0 {
			g.CoroState.delay = delay
			g.CoroState.step = frame
		} else {
			g.GameState = GameStateCoreLoop
		}

	case GameStateCoreLoop:
		g.RenderFrame()
		{
		updateTwice:
			for range 2 {
				var ret Return
				if coro := g.CoroState.coro; coro != nil {
					if g.CoroState.delay > 0 {
						g.CoroState.delay -= 1
					} else if frame, delay := coro(g, g.CoroState.step); frame > 0 {
						g.CoroState.delay = delay
						g.CoroState.step = frame
					} else {
						g.CoroState.coro = nil
						ret = g.CoroState.next(g)
					}
				} else {
					ret = g.UpdateState()
				}

				if ret.coro != nil {
					g.CoroState.coro = ret.coro
					g.CoroState.next = ret.next
					g.CoroState.step = 0
					g.CoroState.delay = 0
				}
				if ret.done {
					break updateTwice
				}
			}
		}
	}
}
