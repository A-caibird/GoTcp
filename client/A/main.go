package main

import (
	. "acaibird.com/clientA/log"
	"acaibird.com/clientA/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net"
	"time"
)

// TODO 客户端如何一直保持链接,以便及时接收服务器的消息
func main() {
	// 连接服务器
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("dial error:", err)
		Logger.Error("服务器连接异常", zap.Error(err))
		return
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			Logger.Error("客户端关闭异常", zap.Error(err))
		}
	}(conn)

	// 发送验证消息
	login := message.TextMsg{
		Type:     "login",
		Sender:   "A",
		Receiver: "server",
		Content:  "client login",
		Time:     time.Now(),
	}

	jsonLogin, _ := json.Marshal(login)
	// 消息长度
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(len(jsonLogin)))
	_, err = conn.Write(buf)       // 发送消息长度
	_, err = conn.Write(jsonLogin) // 发送消息内容

	// A向B发送消息
	msg := message.TextMsg{
		Type:     "text",
		Sender:   "A",
		Receiver: "B",
		Content:  "hello B",
		Time:     time.Now(),
	}

	buf = make([]byte, 8)
	jsonText, _ := json.Marshal(msg)
	binary.BigEndian.PutUint64(buf, uint64(len(jsonText)))
	_, err = conn.Write(buf)
	_, err = conn.Write(jsonText)
}
