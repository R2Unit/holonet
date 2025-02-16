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
)

const (
	token      = "secret123" // must match the core's validToken.
	workerName = "worker1"   // Unique worker name.
)

// dialWebSocket connects to the WebSocket server and performs the handshake.
func dialWebSocket(urlStr string) (net.Conn, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// Append token and worker name as query parameters.
	q := u.Query()
	q.Set("token", token)
	q.Set("name", workerName)
	u.RawQuery = q.Encode()

	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return nil, err
	}
	// Generate a random Sec-WebSocket-Key.
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
	if _, err := conn.Write([]byte(req)); err != nil {
		conn.Close()
		return nil, err
	}
	// Read the server's handshake response.
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
	return conn, nil
}

// readMessage reads a single unfragmented text message from the WebSocket connection.
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
	return string(payload), nil
}

// writeMessage writes a single unfragmented text message to the WebSocket connection.
func writeMessage(conn net.Conn, message string) error {
	payload := []byte(message)
	var header bytes.Buffer
	header.WriteByte(0x81) // FIN + text frame opcode
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

// Task mirrors the structure in the core queue.
type Task struct {
	ID      string            `json:"id"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Files   map[string]string `json:"files,omitempty"`
}

// executeTask writes provided files to a temporary directory (if any),
// sends a "started" status message, executes the task, and then sends a "done" update.
func executeTask(task Task, conn net.Conn) {
	// If there are files, write them to a temporary directory.
	if len(task.Files) > 0 {
		tempDir, err := ioutil.TempDir("", "task-"+task.ID)
		if err != nil {
			log.Printf("Task %s: failed to create temporary directory: %v", task.ID, err)
			return
		}
		defer os.RemoveAll(tempDir) // Clean up after execution.

		// Write each file and update arguments if needed.
		for filename, content := range task.Files {
			filePath := filepath.Join(tempDir, filename)
			if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
				log.Printf("Task %s: failed to write file %s: %v", task.ID, filename, err)
				return
			}
			// Replace file references in Args with the full path.
			for i, arg := range task.Args {
				if arg == filename {
					task.Args[i] = filePath
				}
			}
		}
	}

	// Send a status update: task started.
	startMsg := fmt.Sprintf("Task %s started on worker '%s'", task.ID, workerName)
	if err := writeMessage(conn, startMsg); err != nil {
		log.Println("Error sending start status:", err)
	}

	log.Printf("Executing task %s: %s %v", task.ID, task.Command, task.Args)
	cmd := exec.Command(task.Command, task.Args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Task %s failed: %v\nOutput: %s", task.ID, err, string(output))
		// Optionally, send a failure status update.
		failMsg := fmt.Sprintf("Task %s failed on worker '%s'", task.ID, workerName)
		_ = writeMessage(conn, failMsg)
		return
	}

	log.Printf("Task %s completed successfully:\n%s", task.ID, string(output))
	// Send a completion status update.
	doneMsg := fmt.Sprintf("Task %s completed on worker '%s'", task.ID, workerName)
	if err := writeMessage(conn, doneMsg); err != nil {
		log.Println("Error sending completion status:", err)
	}
}

func main() {
	wsURL := "ws://localhost:8080/ws"
	conn, err := dialWebSocket(wsURL)
	if err != nil {
		log.Fatal("Error connecting to core:", err)
	}
	defer conn.Close()

	log.Println("Worker connected to core via WebSocket")

	// Listen for tasks from the core.
	for {
		msg, err := readMessage(conn)
		if err != nil {
			log.Println("Error reading message (connection lost):", err)
			break // Exit the loop if the connection is lost.
		}
		var task Task
		if err := json.Unmarshal([]byte(msg), &task); err != nil {
			log.Println("Error unmarshalling task:", err)
			continue
		}
		// Execute the task concurrently, passing the connection for status updates.
		go executeTask(task, conn)
	}
	log.Println("Worker: Disconnected from core")
}
