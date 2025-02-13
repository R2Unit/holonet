// core/controller/controller.go
package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/r2unit/talos-core/queue"
	"github.com/r2unit/talos-core/web"
)

const validToken = "secret123"

func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := web.Upgrade(w, r, validToken)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		select {
		case task := <-queue.TaskQueue:
			msg, err := json.Marshal(task)
			if err != nil {
				log.Println("JSON Marshal error:", err)
				continue
			}
			if err := web.WriteMessage(conn, string(msg)); err != nil {
				log.Println("WriteMessage error:", err)
				return
			}
			log.Printf("Sent task %s to worker", task.ID)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
