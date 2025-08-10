package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/blues/jsonata-go"
)

type Input struct {
	JSONInput   map[string]interface{} `json:"json_input"`
	JSONataExpr string                 `json:"jsonata_expr"`
}

type Result struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Error   error       `json:"error"`
	Message string      `json:"message"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var in Input
	if err := decoder.Decode(&in); err != nil {
		encoder.Encode(Result{
			Success: false,
			Error:   err,
			Message: "Invalid JSON",
		})
		return
	}

	expr, err := jsonata.Compile(in.JSONataExpr)
	if err != nil {
		encoder.Encode(Result{
			Success: false,
			Error:   err,
			Message: "Invalid expression",
		})
		return
	}

	resp, err := expr.Eval(in.JSONInput)
	if err != nil {
		encoder.Encode(Result{
			Success: false,
			Error:   err,
			Message: "Evaluation error",
		})
		return
	}

	encoder.Encode(Result{
		Result:  resp,
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
