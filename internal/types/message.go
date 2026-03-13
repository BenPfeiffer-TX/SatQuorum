package types

import "time"

type Message struct {
	ID        string    `json:"id"`
	Payload   string    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}
