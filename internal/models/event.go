package models

type Event interface {
	GetType() EventType
	Handle(EventData) (*Response, error)
}

type EventType uint8

const (
	OnMessageCreate EventType = 1
)

type EventData struct {
	Message *Message
}
