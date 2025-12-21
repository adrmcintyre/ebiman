package game

// A Coro is a Game method that runs one step of a coroutine.
//
// The coroutine is invoked for the first time with step=0.
//
// On return, nextStep is the value to supply for step in the next call to
// the coroutine, and delay is how many frames to wait (if any).
//
// The end of the coroutine's work is signalled when it returns (0, 0).
type Coro func(g *Game, step int) (nextStep int, delay int)

// A CoroState describes the state of a coroutine.
type CoroState struct {
	coro  Coro         // coroutine's invocation method; nil if not active
	step  int          // current step in the coroutine
	delay int          // number of frames to delay before running coro again
	next  Continuation // method to call when the coroutine has finished
}

// A Contination is the signature of a Game method to invoke upon termination
// of a coroutine. The value returned by the continuation indicates the next
// action to take.
type Continuation func(*Game) Return

// A Return allows a state update method to specify what action to take
// after it returns.
type Return struct {
	done bool         // do not advance game state any further during this update period
	coro Coro         // optional coroutine to start
	next Continuation // optional continuation to invoke when coroutine ends
}

// thenStop is a convenience value for a Return indicating that no further
// update processing should occur during this update period.
var thenStop = Return{true, nil, nil}

// thenContinue is a convenience value for a Return indicating that further
// update processing should continue as normal.
var thenContinue = Return{false, nil, nil}

// withCoro is a convenience method which constructs a Return indicating that
// the coroutine coro should be started, and when it terminates, execution
// should continue by invoking next.
func withCoro(coro Coro, next Continuation) Return {
	return Return{false, coro, next}
}
