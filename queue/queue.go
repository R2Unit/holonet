package queue

import "sync"

type Task struct {
	ID           string            `json:"id"`
	Command      string            `json:"command"`
	Args         []string          `json:"args"`
	Files        map[string]string `json:"files,omitempty"`
	Reporter     string            `json:"reporter"`
	Hosts        string            `json:"hosts"`
	TaskTemplate string            `json:"task_template"`
}

var TaskQueue = make(chan Task, 100)

var (
	pendingTasks []Task
	queueMutex   sync.Mutex
)

func Enqueue(task Task) {
	queueMutex.Lock()
	pendingTasks = append(pendingTasks, task)
	queueMutex.Unlock()
	TaskQueue <- task
}

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

func GetPendingTasks() []Task {
	queueMutex.Lock()
	defer queueMutex.Unlock()
	tasksCopy := make([]Task, len(pendingTasks))
	copy(tasksCopy, pendingTasks)
	return tasksCopy
}
