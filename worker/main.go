// worker/main.go
package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const (
	token      = "secret123"              // Must match the core's validToken.
	workerName = "worker1"                // Unique worker name.
	wsURL      = "ws://localhost:8080/ws" // Core's WebSocket URL.
)

var (
	debug               bool
	currentTaskID       string = "none" // holds the current task ID; "none" when idle
	currentTaskTemplate string = ""     // holds the task_template of current task
	currentTaskHosts    string = ""     // holds the hosts of current task
	isRunning           bool
	taskMutex           sync.Mutex
	writeMutex          sync.Mutex
)

func init() {
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}
	if debug {
		log.Println("Debug mode enabled")
	}
}

// Task represents a unit of work to be executed.
type Task struct {
	ID           string            `json:"id"`
	Command      string            `json:"command"`
	Args         []string          `json:"args"`
	Files        map[string]string `json:"files,omitempty"`
	Reporter     string            `json:"reporter"`      // provided via API task creation
	Hosts        string            `json:"hosts"`         // provided via API task creation
	TaskTemplate string            `json:"task_template"` // e.g. "ansible"
}

// WorkerStatus is sent from the worker to the core.
type WorkerStatus struct {
	Worker       string `json:"worker"`
	TaskID       string `json:"task_id,omitempty"`
	Status       string `json:"status"`                  // "running" or "idle"
	Hosts        string `json:"hosts,omitempty"`         // included when running
	TaskTemplate string `json:"task_template,omitempty"` // included when running
	Reporter     string `json:"reporter,omitempty"`      // included when running (set to "automation")
}

// safeWriteMessage locks and writes a message.
func safeWriteMessage(conn net.Conn, message string) error {
	writeMutex.Lock()
	defer writeMutex.Unlock()
	return writeMessage(conn, message)
}

// sendStatus marshals and sends a status update.
func sendStatus(conn net.Conn, status WorkerStatus) error {
	data, err := json.Marshal(status)
	if err != nil {
		log.Println("Error marshaling status:", err)
		return err
	}
	if debug {
		log.Printf("[DEBUG] Sending status: %s", string(data))
	}
	return safeWriteMessage(conn, string(data))
}

// dialWebSocket performs the websocket handshake.
func dialWebSocket(urlStr string) (net.Conn, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("token", token)
	q.Set("name", workerName)
	u.RawQuery = q.Encode()

	if debug {
		log.Printf("[DEBUG] Attempting connection to %s", u.String())
	}
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return nil, err
	}
	keyBytes := make([]byte, 16)
	if _, err := rand.Read(keyBytes); err != nil {
		conn.Close()
		return nil, err
	}
	key := base64.StdEncoding.EncodeToString(keyBytes)
	req := fmt.Sprintf("GET %s HTTP/1.1\r\n", u.RequestURI()) +
		fmt.Sprintf("Host: %s\r\n", u.Host) +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: " + key + "\r\n" +
		"Sec-WebSocket-Version: 13\r\n\r\n"
	if debug {
		log.Printf("[DEBUG] Sending handshake request:\n%s", req)
	}
	if _, err := conn.Write([]byte(req)); err != nil {
		conn.Close()
		return nil, err
	}
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, &http.Request{Method: "GET"})
	if err != nil {
		conn.Close()
		return nil, err
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		conn.Close()
		return nil, errors.New("failed to upgrade to websocket")
	}
	if debug {
		log.Printf("[DEBUG] Connection upgraded, response: %s", resp.Status)
	}
	return conn, nil
}

// readMessage reads a single websocket text message.
func readMessage(conn net.Conn) (string, error) {
	var header [2]byte
	if _, err := io.ReadFull(conn, header[:2]); err != nil {
		return "", err
	}
	fin := header[0] & 0x80
	opcode := header[0] & 0x0F
	if fin == 0 || opcode != 1 {
		return "", errors.New("only single-frame text messages supported")
	}
	masked := header[1] & 0x80
	payloadLen := int(header[1] & 0x7F)
	if payloadLen == 126 {
		var ext [2]byte
		if _, err := io.ReadFull(conn, ext[:2]); err != nil {
			return "", err
		}
		payloadLen = int(binary.BigEndian.Uint16(ext[:2]))
	} else if payloadLen == 127 {
		var ext [8]byte
		if _, err := io.ReadFull(conn, ext[:8]); err != nil {
			return "", err
		}
		payloadLen = int(binary.BigEndian.Uint64(ext[:8]))
	}
	var maskingKey [4]byte
	if masked != 0 {
		if _, err := io.ReadFull(conn, maskingKey[:4]); err != nil {
			return "", err
		}
	}
	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return "", err
	}
	if masked != 0 {
		for i := 0; i < payloadLen; i++ {
			payload[i] ^= maskingKey[i%4]
		}
	}
	if debug {
		log.Printf("[DEBUG] Received message: %s", string(payload))
	}
	return string(payload), nil
}

// writeMessage writes a websocket text message.
func writeMessage(conn net.Conn, message string) error {
	if debug {
		log.Printf("[DEBUG] Sending message: %s", message)
	}
	payload := []byte(message)
	var header bytes.Buffer
	header.WriteByte(0x81) // FIN + text opcode
	length := len(payload)
	if length < 126 {
		header.WriteByte(byte(length))
	} else if length <= 65535 {
		header.WriteByte(126)
		var ext [2]byte
		binary.BigEndian.PutUint16(ext[:], uint16(length))
		header.Write(ext[:])
	} else {
		header.WriteByte(127)
		var ext [8]byte
		binary.BigEndian.PutUint64(ext[:], uint64(length))
		header.Write(ext[:])
	}
	if _, err := conn.Write(header.Bytes()); err != nil {
		return err
	}
	_, err := conn.Write(payload)
	return err
}

// heartbeat sends status updates every second.
func heartbeat(conn net.Conn, done chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			taskMutex.Lock()
			tID := currentTaskID
			running := isRunning
			tTemplate := currentTaskTemplate
			tHosts := currentTaskHosts
			taskMutex.Unlock()
			var ws WorkerStatus
			if running {
				ws = WorkerStatus{
					Worker:       workerName,
					TaskID:       tID,
					Status:       "running",
					Hosts:        tHosts,
					Reporter:     "automation",
					TaskTemplate: tTemplate,
				}
			} else {
				ws = WorkerStatus{
					Worker: workerName,
					TaskID: tID,
					Status: "idle",
				}
			}
			if err := sendStatus(conn, ws); err != nil {
				log.Println("[Heartbeat] Error sending status, exiting heartbeat")
				return
			}
		}
	}
}

// executeTask writes any files, runs galaxy installation if a requirements.yml exists,
// then runs the main command (ansible-playbook), and sends status updates.
func executeTask(task Task, conn net.Conn) {
	if debug {
		log.Printf("[DEBUG] Starting execution of task: %+v", task)
	}
	var tempDir string
	if len(task.Files) > 0 {
		var err error
		tempDir, err = ioutil.TempDir("", "task-"+task.ID)
		if err != nil {
			log.Printf("Task %s: failed to create temporary directory: %v", task.ID, err)
			return
		}
		if debug {
			log.Printf("[DEBUG] Temporary directory created: %s", tempDir)
		}
		defer os.RemoveAll(tempDir)
		for filename, content := range task.Files {
			filePath := filepath.Join(tempDir, filename)
			if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
				log.Printf("Task %s: failed to write file %s: %v", task.ID, filename, err)
				return
			}
			if debug {
				log.Printf("[DEBUG] Wrote file %s (length %d)", filePath, len(content))
			}
			for i, arg := range task.Args {
				if arg == filename {
					if debug {
						log.Printf("[DEBUG] Replacing argument %s with %s", arg, filePath)
					}
					task.Args[i] = filePath
				}
			}
		}
	}

	// Update global task info.
	taskMutex.Lock()
	currentTaskID = task.ID
	currentTaskTemplate = task.TaskTemplate
	currentTaskHosts = task.Hosts
	isRunning = true
	taskMutex.Unlock()

	// Before running the playbook, check for requirements.yml and run galaxy install.
	if _, ok := task.Files["requirements.yml"]; ok && tempDir != "" {
		reqPath := filepath.Join(tempDir, "requirements.yml")
		if debug {
			log.Printf("[DEBUG] Found requirements.yml at %s, running ansible-galaxy install", reqPath)
		}
		cmdGalaxy := exec.Command("ansible-galaxy", "install", "-r", reqPath, "--force")
		galaxyOutput, galaxyErr := cmdGalaxy.CombinedOutput()
		if galaxyErr != nil {
			log.Printf("Task %s: ansible-galaxy install failed: %v\nOutput: %s", task.ID, galaxyErr, string(galaxyOutput))
			taskMutex.Lock()
			isRunning = false
			currentTaskID = "none"
			currentTaskTemplate = ""
			currentTaskHosts = ""
			taskMutex.Unlock()
			sendStatus(conn, WorkerStatus{
				Worker: workerName,
				TaskID: "none",
				Status: "idle",
			})
			return
		}
		if debug {
			log.Printf("[DEBUG] ansible-galaxy install output: %s", string(galaxyOutput))
		}
	}

	// Send running status update with extra fields.
	if err := sendStatus(conn, WorkerStatus{
		Worker:       workerName,
		TaskID:       task.ID,
		Status:       "running",
		Hosts:        task.Hosts,
		Reporter:     "automation",
		TaskTemplate: task.TaskTemplate,
	}); err != nil {
		log.Println("Error sending running status:", err)
	}

	log.Printf("Executing task %s: %s %v", task.ID, task.Command, task.Args)
	cmd := exec.Command(task.Command, task.Args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Task %s failed: %v\nOutput: %s", task.ID, err, string(output))
		taskMutex.Lock()
		isRunning = false
		currentTaskID = "none"
		currentTaskTemplate = ""
		currentTaskHosts = ""
		taskMutex.Unlock()
		sendStatus(conn, WorkerStatus{
			Worker: workerName,
			TaskID: "none",
			Status: "idle",
		})
		return
	}
	log.Printf("Task %s completed successfully:\n%s", task.ID, string(output))
	taskMutex.Lock()
	isRunning = false
	currentTaskID = "none"
	currentTaskTemplate = ""
	currentTaskHosts = ""
	taskMutex.Unlock()
	if err := sendStatus(conn, WorkerStatus{
		Worker: workerName,
		TaskID: "none",
		Status: "idle",
	}); err != nil {
		log.Println("Error sending idle status after completion:", err)
	}
}

func main() {
	for {
		conn, err := dialWebSocket(wsURL)
		if err != nil {
			log.Println("Error connecting to core:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Worker connected to core via WebSocket")

		heartbeatDone := make(chan struct{})
		go heartbeat(conn, heartbeatDone)

		for {
			msg, err := readMessage(conn)
			if err != nil {
				log.Println("Error reading message (connection lost):", err)
				break
			}
			var task Task
			if err := json.Unmarshal([]byte(msg), &task); err != nil {
				log.Println("Error unmarshalling task:", err)
				continue
			}
			go executeTask(task, conn)
		}
		close(heartbeatDone)
		conn.Close()
		log.Println("Worker: Disconnected from core, retrying in 5 seconds")
		time.Sleep(5 * time.Second)
	}
}
