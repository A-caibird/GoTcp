package message

import "time"

type MSG interface {
	getSender() string
	getReceiver() string
}
type TextMsg struct {
	Type     string    `json:"type"`
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}

func (t TextMsg) GetSender() string {
	return t.Sender
}

func (t TextMsg) GetReceiver() string {
	return t.Receiver
}

func (t TextMsg) GetContent() string {
	return t.Content
}
