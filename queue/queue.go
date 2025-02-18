// core/queue/queue.go
package queue

import "sync"

// Task represents a unit of work that a worker will execute.
type Task struct {
	ID           string            `json:"id"`
	Command      string            `json:"command"`
	Args         []string          `json:"args"`
	Files        map[string]string `json:"files,omitempty"`
	Reporter     string            `json:"reporter"`      // who requested this task
	Hosts        string            `json:"hosts"`         // hosts as given by the API caller
	TaskTemplate string            `json:"task_template"` // type of task (e.g. "ansible")
}

var TaskQueue = make(chan Task, 100)

var (
	pendingTasks []Task
	queueMutex   sync.Mutex
)

// Enqueue adds a new task to the queue.
func Enqueue(task Task) {
	queueMutex.Lock()
	pendingTasks = append(pendingTasks, task)
	queueMutex.Unlock()
	TaskQueue <- task
}

// Dequeue removes a task from the pendingTasks slice by ID.
func Dequeue(taskID string) {
	queueMutex.Lock()
	defer queueMutex.Unlock()
	for i, t := range pendingTasks {
		if t.ID == taskID {
			pendingTasks = append(pendingTasks[:i], pendingTasks[i+1:]...)
			break
		}
	}
}

// GetPendingTasks returns a copy of the pending tasks.
func GetPendingTasks() []Task {
	queueMutex.Lock()
	defer queueMutex.Unlock()
	tasksCopy := make([]Task, len(pendingTasks))
	copy(tasksCopy, pendingTasks)
	return tasksCopy
}
