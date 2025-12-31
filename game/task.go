package game

// A TaskId identifies a task to perform before regular state updates.
type TaskId int

// The possible tasks (currently the only one!)
const (
	TaskGhostReturn TaskId = iota // mark a ghost as having returned
)

// A Task is a piece of work to perform before the regular state update.
type Task struct {
	Id    TaskId // identifies the type of task
	Param int    // simple parameter storage
}

// AddTask adds a new task to the queue.
func (g *Game) AddTask(id TaskId, param int) {
	g.TaskQueue = append(g.TaskQueue, Task{id, param})
}

// RunTaskQueue executes all tasks pending in the queue.
func (g *Game) RunTaskQueue() {
	for len(g.TaskQueue) > 0 {
		switch task := g.TaskQueue[0]; task.Id {
		case TaskGhostReturn:
			// TODO - the only task - incorporate the delay processing here too?
			g.GhostReturn(task.Param)
		}
		g.TaskQueue = g.TaskQueue[1:]
	}
}
