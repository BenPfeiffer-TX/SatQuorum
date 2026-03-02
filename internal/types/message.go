package types

type message struct {
	ID        string `json:"id"`
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
}
