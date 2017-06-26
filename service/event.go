package service

type EventType string

const (
	TextEvent EventType = "text"
)

// Event structure
type Event struct {
	Type EventType
	From string
	To   string
	Data []byte
	Meta map[string]interface{}
}
