package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/blues/jsonata-go"
)

type Input struct {
	JSONInput   map[string]interface{} `json:"json_input"`
	JSONataExpr string                 `json:"jsonata_expr"`
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Error reading stdin:", err)
	}

	var in Input
	if err := json.Unmarshal(data, &in); err != nil {
		log.Fatal("Invalid JSON:", err)
	}

	if in.JSONataExpr == "" {
		log.Fatal("Missing jsonata_expr in input")
	}

	expr, err := jsonata.Compile(in.JSONataExpr)
	if err != nil {
		log.Fatal("Invalid JSONata expression:", err)
	}

	resp, err := expr.Eval(in.JSONInput)
	if err != nil {
		log.Fatal("Evaluation error:", err)
	}

	out, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("Error encoding output:", err)
	}

	os.Stdout.Write(out)
}
