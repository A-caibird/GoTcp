package message

import (
	"acaibird.com/server/mysql"
	"errors"
	"time"
)

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

func (t TextMsg) GetContent() interface{} {
	return t.Content
}

func (t TextMsg) WriteToDB() (err error) {
	db, err := mysqlDB.InitDB()
	if err != nil {
		return errors.New("数据库连接异常")
	}
	defer func() {
		err = db.Close()
		if err != nil {
			return // 如果数据库关闭失败,也返回了对应的err
		}
	}()

	// 准备 SQL 语句
	stmt, err := db.Prepare("INSERT INTO text_msgs (type, sender, receiver, content, time) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("数据库准备 SQL 语句异常")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	// 执行 SQL 语句
	_, err = stmt.Exec(t.Type, t.Sender, t.Receiver, t.Content, t.Time)
	if err != nil {
		return errors.New("数据库执行 SQL 语句异常")
	}
	return err // 如果数据库关闭失败,也返回了对应的err
}

func (t TextMsg) GetTime() time.Time {
	return t.Time
}

type FileMsg struct {
	Type        string    `json:"type"`
	Sender      string    `json:"sender"`
	Receiver    string    `json:"receiver"`
	FileName    string    `json:"file_name"`
	FileContent int64     `json:"file_content"`
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
