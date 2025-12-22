package game

import "time"

// A CoroMethod is a Game method that runs one step of a coroutine.
//
// The coro argument maintains the current state of the coroutine,
// and provides methods for indicating how to continue upon next
// invocation.
type CoroMethod func(g *Game, coro *Coro) *Coro

// A Coro describes the state of a coroutine.
type Coro struct {
	method    CoroMethod   // coroutine's invocation method
	step      int          // current step in the coroutine
	waitUntil time.Time    // do not invoke coro again until after this time
	next      Continuation // method to call when the coroutine has finished
}

// Step returns the current step to execute.
func (coro *Coro) Step() int {
	return coro.step
}

// WaitNext causes coro's next invocation to be delayed by ms milliseconds,
// with an incremented step.
func (coro *Coro) WaitNext(ms int) *Coro {
	return coro.Wait(ms).Next()
}

// Next causes coro's next invocation to be with an incremented step.
func (coro *Coro) Next() *Coro {
	coro.step += 1
	return coro
}

// Wait causes coro's next invocation to be delayed by ms milliseconds,
// retaining the current step.
func (coro *Coro) Wait(ms int) *Coro {
	coro.waitUntil = time.Now().Add(time.Duration(ms) * time.Millisecond)
	return coro
}

// Stop causes the coro to terminate, and the next continuation to be called.
func (coro *Coro) Stop() *Coro {
	return nil
}

// invoke invokes coro's method if the wait time has expired.
func (coro *Coro) invoke(g *Game) *Coro {
	if time.Now().Before(coro.waitUntil) {
		return coro
	}
	return coro.method(g, coro)
}

// A Contination is the signature of a Game method to invoke upon termination
// of a coroutine. The continuation's return value indicates the next action
// to take.
type Continuation func(*Game) Return

// A Return allows a state update method to specify what action to take
// after it returns.
type Return struct {
	done bool  // do not advance game state any further during this update period
	coro *Coro // optional coroutine
}

// thenStop is a convenience value for a Return indicating that no further
// update processing should occur during this update period.
var thenStop = Return{done: true}

// thenContinue is a convenience value for a Return indicating that further
// update processing should continue as normal.
var thenContinue = Return{done: false}

// withCoro is a convenience method which constructs a Return indicating that
// a coroutine for method should be started, and when it terminates, execution
// should continue by invoking next.
func withCoro(method CoroMethod, next Continuation) Return {
	return Return{
		done: false,
		coro: &Coro{
			method: method,
			next:   next,
		},
	}
}
