package message

import "time"

type MessageTypeTip struct {
	Index int64
	Name  string
}

type MSG interface {
	GetSender() string
	GetReceiver() string
	GetTime() time.Time
	GetContent() interface{}
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
func (t TextMsg) GetTime() time.Time {
	return t.Time
}

type FileMsg struct {
	Type        string    `json:"type"`
	Sender      string    `json:"sender"`
	Receiver    string    `json:"receiver"`
	FileName    string    `json:"file_name"`
	FileContent []byte    `json:"file_content"`
	Time        time.Time `json:"time"`
}

func (f FileMsg) GetSender() string {
	return f.Sender
}
func (f FileMsg) GetReceiver() string {
	return f.Receiver
}
func (f FileMsg) GetContent() interface{} {
	return f.FileContent
}
func (f FileMsg) GetTime() time.Time {
	return f.Time
}
