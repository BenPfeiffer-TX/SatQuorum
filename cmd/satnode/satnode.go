package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	var msg map[string]interface{} //replace with actual message struct
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := map[string]string{"status": "received"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to encode response:", err)
	}
}
