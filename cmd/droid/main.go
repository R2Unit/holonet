package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func generateWebSocketKey() (string, error) {
	keyBytes := make([]byte, 16)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(keyBytes), nil
}

// readFrame reads a single WebSocket frame from conn.
// This simplified version assumes that the frame is not masked and contains a small payload.
func readFrame(conn net.Conn) (opcode byte, payload []byte, err error) {
	header := make([]byte, 2)
	if _, err = io.ReadFull(conn, header); err != nil {
		return
	}

	// fin and opcode extraction.
	opcode = header[0] & 0x0F
	payloadLen := int(header[1] & 0x7F)

	// If the payload length is extended into the next bytes.
	if payloadLen == 126 {
		ext := make([]byte, 2)
		if _, err = io.ReadFull(conn, ext); err != nil {
			return
		}
		payloadLen = int(binary.BigEndian.Uint16(ext))
	} else if payloadLen == 127 {
		ext := make([]byte, 8)
		if _, err = io.ReadFull(conn, ext); err != nil {
			return
		}
		payloadLen = int(binary.BigEndian.Uint64(ext))
	}

	payload = make([]byte, payloadLen)
	_, err = io.ReadFull(conn, payload)
	return
}

func main() {
	// Use command-line flags for the WebSocket server's host:port and token.
	url := flag.String("url", "localhost:8080", "WebSocket server host:port (e.g., localhost:8080)")
	token := flag.String("token", "test", "Token for connecting to the WebSocket")
	flag.Parse()

	// Generate a random Sec-WebSocket-Key.
	secWebSocketKey, err := generateWebSocketKey()
	if err != nil {
		log.Fatalf("error generating Sec-WebSocket-Key: %v", err)
	}

	// Prepare the HTTP/WebSocket handshake request.
	request := fmt.Sprintf("GET /?token=%s HTTP/1.1\r\n"+
		"Host: %s\r\n"+
		"Upgrade: websocket\r\n"+
		"Connection: Upgrade\r\n"+
		"Sec-WebSocket-Key: %s\r\n"+
		"Sec-WebSocket-Version: 13\r\n\r\n", *token, *url, secWebSocketKey)

	// Connect to the WebSocket server.
	conn, err := net.Dial("tcp", *url)
	if err != nil {
		log.Fatalf("error dialing %s: %v", *url, err)
	}
	defer conn.Close()

	// Send the handshake request.
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatalf("error writing handshake request: %v", err)
	}

	// Read and print the HTTP response:
	reader := bufio.NewReader(conn)

	// Read the status line.
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading status line: %v", err)
	}
	if !strings.Contains(statusLine, "101") {
		log.Fatalf("handshake failed, expected status 101 but got: %s", statusLine)
	}
	fmt.Printf("Handshake response: %s", statusLine)

	// Read and discard the rest of the headers.
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading header: %v", err)
		}
		if line == "\r\n" {
			break
		}
	}

	// Pause briefly to ensure the server sends its first frame.
	time.Sleep(50 * time.Millisecond)

	// Read a WebSocket frame from the server.
	opcode, payload, err := readFrame(conn)
	if err != nil {
		log.Fatalf("error reading WebSocket frame: %v", err)
	}

	// Currently, opcode 1 represents a text frame.
	if opcode != 1 {
		log.Printf("expected text frame (opcode 1), got opcode %d", opcode)
	} else {
		fmt.Printf("Received text frame: %s\n", string(payload))
	}
}
