// core/web/websocket.go
package web

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
)

func ComputeAcceptKey(key string) string {
	const magic = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.New()
	h.Write([]byte(key + magic))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Upgrade(w http.ResponseWriter, r *http.Request, validToken string) (net.Conn, error) {
	token := r.URL.Query().Get("token")
	name := r.URL.Query().Get("name")
	if token != validToken || name == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil, errors.New("invalid token or missing worker name")
	}
	log.Printf("Worker '%s' attempting to connect", name)

	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return nil, errors.New("missing Sec-WebSocket-Key")
	}
	acceptKey := ComputeAcceptKey(key)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("server does not support hijacking")
	}
	conn, bufrw, err := hijacker.Hijack()
	if err != nil {
		return nil, err
	}

	response := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + acceptKey + "\r\n\r\n"
	if _, err := bufrw.WriteString(response); err != nil {
		conn.Close()
		return nil, err
	}
	if err := bufrw.Flush(); err != nil {
		conn.Close()
		return nil, err
	}
	log.Printf("Worker '%s' connected successfully", name)
	return conn, nil
}

func ReadMessage(conn net.Conn) (string, error) {
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

func WriteMessage(conn net.Conn, message string) error {
	payload := []byte(message)
	var header bytes.Buffer
	header.WriteByte(0x81)
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
