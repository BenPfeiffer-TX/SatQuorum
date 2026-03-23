package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/BenPfeiffer-TX/SatQuorum/internal/types"
)

/*
This is the main service that the 'satellites' will be running in their docker containers
it will listen on a port to receive some structured message (probably json)
it will do a calculation based on the received message
it will send a response
*/
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	http.HandleFunc("/", fooHandler)

	log.Printf("Server starting on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	//var msg map[string]interface{} //replace with actual message struct
	var msg types.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Printf("failed to decode message")
		return
	}

	log.Printf("received a message: ", msg.ID, msg.Payload, msg.Timestamp)
	response := map[string]string{"status": "received"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response:", err)
	}
}
