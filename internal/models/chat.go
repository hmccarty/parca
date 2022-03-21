package models

type ChatClient interface{}

type ChatMessage struct {
	IsDM     bool
	ID       string
	Content  string
	Reaction string
	Values   []string
}
