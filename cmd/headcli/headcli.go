package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BenPfeiffer-TX/SatQuorum/internal/types"
)

/*
control software intended to send messages to satellite nodes
*/

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	println(addr)
	msg := types.Message{ID: "doctor", Payload: "this is a test, did you pass?", Timestamp: time.Now()}
	sendMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost"+addr, bytes.NewReader(sendMsg))
	if err != nil {
		log.Println(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()
}
