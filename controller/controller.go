package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/r2unit/colours"
	"github.com/r2unit/holonet-core/queue"
	"github.com/r2unit/holonet-core/web"
)

const validToken = "insecure"

type WorkerInfo struct {
	Name        string `json:"name"`
	CurrentTask string `json:"current_task"`
	Health      string `json:"health"`
}

var (
	workerRegistry      = make(map[string]*WorkerInfo)
	workerRegistryMutex sync.Mutex

	db *sql.DB
)

func SetDB(d *sql.DB) {
	db = d
}

func RegisterWorker(name string) {
	workerRegistryMutex.Lock()
	defer workerRegistryMutex.Unlock()
	workerRegistry[name] = &WorkerInfo{
		Name:        name,
		CurrentTask: "none",
		Health:      "connected",
	}
}

func UpdateWorkerTask(name, taskID string) {
	workerRegistryMutex.Lock()
	defer workerRegistryMutex.Unlock()
	if w, ok := workerRegistry[name]; ok {
		if taskID == "" || taskID == "none" {
			w.CurrentTask = "none"
		} else {
			w.CurrentTask = taskID
		}
	}
}

func SetWorkerHealth(name, health string) {
	workerRegistryMutex.Lock()
	defer workerRegistryMutex.Unlock()
	if w, ok := workerRegistry[name]; ok {
		w.Health = health
	}
}

func GetWorkers() []WorkerInfo {
	workerRegistryMutex.Lock()
	defer workerRegistryMutex.Unlock()
	workers := []WorkerInfo{}
	for _, w := range workerRegistry {
		workers = append(workers, *w)
	}
	return workers
}

func InsertWorkerLog(worker, task, status, hosts, taskTemplate, reporter string) {
	if db == nil {
		log.Println("DB not initialized, skipping log insertion")
		return
	}
	_, err := db.Exec(
		"INSERT INTO workers(worker, task, status, hosts, task_template, reporter) VALUES($1, $2, $3, $4, $5, $6)",
		worker, task, status, hosts, taskTemplate, reporter,
	)
	if err != nil {
		log.Println("Error inserting worker log:", err)
	}
}

type WorkerStatus struct {
	Worker       string `json:"worker"`
	TaskID       string `json:"task_id,omitempty"`
	Status       string `json:"status"`
	Hosts        string `json:"hosts,omitempty"`
	TaskTemplate string `json:"task_template,omitempty"`
	Reporter     string `json:"reporter,omitempty"`
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	workerName := r.URL.Query().Get("name")
	if workerName == "" {
		http.Error(w, "Missing worker name", http.StatusBadRequest)
		return
	}

	conn, err := web.Upgrade(w, r, validToken)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	RegisterWorker(workerName)
	SetWorkerHealth(workerName, "connected")
	log.Printf(colours.Info("Worker '%s' connected"), workerName)

	defer func() {
		conn.Close()
		SetWorkerHealth(workerName, "error")
		log.Printf("Worker '%s' disconnected", workerName)
		InsertWorkerLog(workerName, "", "disconnected", "", "", "")
	}()

	go func() {
		for {
			msg, err := web.ReadMessage(conn)
			if err != nil {
				log.Printf(colours.Info("Worker '%s' disconnected"), workerName)
				SetWorkerHealth(workerName, "error")
				InsertWorkerLog(workerName, "", "disconnected", "", "", "")
				return
			}
			var status WorkerStatus
			if err := json.Unmarshal([]byte(msg), &status); err != nil {
				log.Printf("Non-JSON status from worker '%s': %s", workerName, msg)
			} else {
				if !(status.Status == "idle" && (status.TaskID == "" || status.TaskID == "none")) {
					log.Printf(colours.Info("Status update from worker '%s': task: %s, status: %s, hosts: %s, template: %s, reporter: %s"),
						workerName, status.TaskID, status.Status, status.Hosts, status.TaskTemplate, status.Reporter)
					InsertWorkerLog(status.Worker, status.TaskID, status.Status, status.Hosts, status.TaskTemplate, status.Reporter)
				}
				UpdateWorkerTask(status.Worker, status.TaskID)
				SetWorkerHealth(status.Worker, "connected")
			}
		}
	}()

	for {
		select {
		case task := <-queue.TaskQueue:
			queue.Dequeue(task.ID)
			msg, err := json.Marshal(task)
			if err != nil {
				log.Println("JSON Marshal error:", err)
				continue
			}
			if err := web.WriteMessage(conn, string(msg)); err != nil {
				log.Println("WriteMessage error:", err)
				return
			}
			log.Printf(colours.Info("Sent task %s to worker '%s'"), task.ID, workerName)
			InsertWorkerLog(workerName, task.ID, "sent", task.Hosts, task.TaskTemplate, task.Reporter)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
