package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/blues/jsonata-go"
)

type TransformRequest struct {
	Data       interface{} `json:"data"`
	Expression string      `json:"expression"`
}

type TransformResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   error       `json:"error"`
	Message string      `json:"message"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var req TransformRequest
	if err := decoder.Decode(&req); err != nil {
		encoder.Encode(TransformResponse{
			Success: false,
			Error:   err,
			Message: "Invalid JSON",
		})
		return
	}

	expr, err := jsonata.Compile(req.Expression)
	if err != nil {
		encoder.Encode(TransformResponse{
			Success: false,
			Error:   err,
			Message: "Invalid expression",
		})
		return
	}

	resp, err := expr.Eval(req.Data)
	if err != nil {
		encoder.Encode(TransformResponse{
			Success: false,
			Error:   err,
			Message: "Evaluation error",
		})
		return
	}

	encoder.Encode(TransformResponse{
		Data:    resp,
		Success: true,
		Error:   nil,
	})
}

func main() {
	socketPath := "/tmp/jsonata.sock"

	// Remove socket if it exists
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	defer listener.Close()

	log.Println("JSONata Unix socket server started:", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
