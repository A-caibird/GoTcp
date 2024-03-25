package main

import (
	"acaibird.com/clientA/message"
	. "acaibird.com/clientB/log"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"net"
	"time"
)

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
		Sender:   "B",
		Receiver: "server",
		Content:  "client login",
		Time:     time.Now(),
	}
	jsonLogin, _ := json.Marshal(login)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(len(jsonLogin)))
	_, err = conn.Write(buf)       // 发送消息长度
	_, err = conn.Write(jsonLogin) // 发送消息内容

	// 接收消息
	buf = make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		Logger.Error("客户端接受消息长度异常", zap.Error(err))
		return
	}
	lens := binary.BigEndian.Uint64(buf[:n])

	// 接收消息内容
	buf = make([]byte, lens)
	n, err = conn.Read(buf) // 收到消息
	var msgReceive message.TextMsg
	err = json.Unmarshal(buf[:n], &msgReceive)

	color.Yellow("收到消息:%#v", msgReceive)
}
