package message

import (
	mysqlDB "acaibird.com/server/mysql"
	"errors"
	"time"
)

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

func (t TextMsg) WriteToDB(errCh chan error) {
	db, err := mysqlDB.InitDB()
	if err != nil {
		errCh <- errors.New("数据库连接异常")
		return
	}
	defer func() {
		err = db.Close()
		if err != nil {
			return
		}
	}()

	// 准备 SQL 语句
	stmt, err := db.Prepare("INSERT INTO text_msgs (type, sender, receiver, content, time) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		errCh <- errors.New("数据库准备 SQL 语句异常")
		return
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			return
		}

	}()

	// 执行 SQL 语句
	_, err = stmt.Exec(t.Type, t.Sender, t.Receiver, t.Content, t.Time)
	if err != nil {
		errCh <- errors.New("数据库执行 SQL 语句异常")
		return
	}
	errCh <- nil
}
