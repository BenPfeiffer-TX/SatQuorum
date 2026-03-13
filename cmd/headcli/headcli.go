package main

import (
	"encoding/json"
	"log"
	"net"
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

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	_, err = conn.Write(sendMsg)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
