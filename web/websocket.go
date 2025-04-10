package web

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

// magicKey is a GUID used in WebSocket handshake to compute the Sec-WebSocket-Accept header value.
const magicKey = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// webSocketHandshake performs the WebSocket handshake, validating the request and returning a hijacked network connection.
func webSocketHandshake(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("method not allowed")
	}

	if strings.ToLower(r.Header.Get("Upgrade")) != "websocket" ||
		!strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return nil, fmt.Errorf("bad request")
	}

	secWebSocketKey := r.Header.Get("Sec-WebSocket-Key")
	if secWebSocketKey == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return nil, fmt.Errorf("missing Sec-WebSocket-Key")
	}

	h := sha1.New()
	h.Write([]byte(secWebSocketKey + magicKey))
	acceptKey := base64.StdEncoding.EncodeToString(h.Sum(nil))

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, fmt.Errorf("hijacking not supported")
	}
	conn, buf, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}

	response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\n"+
		"Upgrade: websocket\r\n"+
		"Connection: Upgrade\r\n"+
		"Sec-WebSocket-Accept: %s\r\n\r\n", acceptKey)
	if _, err := buf.WriteString(response); err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	if err := buf.Flush(); err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return conn, nil
}

// writeFrame writes a WebSocket frame with the specified opcode and payload to the given connection.
// It constructs the frame header based on payload length and writes it followed by the payload.
// Returns an error if writing to the connection fails.
func writeFrame(conn net.Conn, opcode byte, payload []byte) error {
	fin := byte(0x80)
	header := []byte{fin | opcode}

	payloadLen := len(payload)
	if payloadLen < 126 {
		header = append(header, byte(payloadLen))
	} else if payloadLen <= 65535 {
		header = append(header, 126, byte(payloadLen>>8), byte(payloadLen))
	} else {
		header = append(header, 127)
		for i := 7; i >= 0; i-- {
			header = append(header, byte(payloadLen>>(8*i)))
		}
	}

	if _, err := conn.Write(header); err != nil {
		return err
	}
	_, err := conn.Write(payload)
	return err
}

// readFrame reads a WebSocket frame from the provided connection, parsing its opcode, payload, and handling masking.
func readFrame(conn net.Conn) (opcode byte, payload []byte, err error) {
	header := make([]byte, 2)
	if _, err = io.ReadFull(conn, header); err != nil {
		return
	}

	fin := header[0] & 0x80
	opcode = header[0] & 0x0F
	maskFlag := header[1] & 0x80
	payloadLen := int(header[1] & 0x7F)

	if payloadLen == 126 {
		ext := make([]byte, 2)
		if _, err = io.ReadFull(conn, ext); err != nil {
			return
		}
		payloadLen = int(ext[0])<<8 | int(ext[1])
	} else if payloadLen == 127 {
		ext := make([]byte, 8)
		if _, err = io.ReadFull(conn, ext); err != nil {
			return
		}
		payloadLen = int(ext[7])
	}

	var maskKey []byte
	if maskFlag != 0 {
		maskKey = make([]byte, 4)
		if _, err = io.ReadFull(conn, maskKey); err != nil {
			return
		}
	}

	payload = make([]byte, payloadLen)
	if _, err = io.ReadFull(conn, payload); err != nil {
		return
	}

	if maskFlag != 0 {
		for i := 0; i < payloadLen; i++ {
			payload[i] ^= maskKey[i%4]
		}
	}

	if fin == 0x00 {
		err = fmt.Errorf("fragmented frames not supported")
	}
	return
}

// HandleWebSocket handles a WebSocket connection by initiating the handshake, managing frames, and responding to events.
// It supports text messages, ping-pong for heartbeat, and closes the connection on receiving a close frame or errors.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := webSocketHandshake(w, r)
	if err != nil {
		log.Printf("Handshake error: %v", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	stopCh := make(chan struct{})
	go StartHeartbeat(conn, stopCh)

	for {
		op, payload, err := readFrame(conn)
		if err != nil {
			log.Printf("Error reading frame: %v", err)
			close(stopCh)
			return
		}

		switch op {
		case 0x1: // Text message
			log.Printf("Received text message: %s", string(payload))
			if err := writeFrame(conn, 0x1, payload); err != nil {
				log.Printf("Error sending text message: %v", err)
				close(stopCh)
				return
			}
		case 0x9: // Ping frame; reply with pong.
			log.Println("Received ping; replying with pong")
			if err := writeFrame(conn, 0xA, payload); err != nil {
				log.Printf("Error sending pong: %v", err)
				close(stopCh)
				return
			}
		case 0xA: // Pong frame; log acknowledgement.
			log.Println("Received pong (heartbeat acknowledgement)")
		case 0x8: // Close frame; terminate the connection.
			log.Println("Received close frame; closing connection")
			close(stopCh)
			return
		default:
			log.Printf("Received unsupported opcode: %d", op)
		}
	}
}
