package queue

type Task struct {
	ID      string            `json:"id"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Files   map[string]string `json:"files,omitempty"`
}

var TaskQueue = make(chan Task, 100)

func Enqueue(task Task) {
	TaskQueue <- task
}
