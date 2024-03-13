package main

import (
	. "acaibird.com/meassege"
	. "acaibird.com/zaplog"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"net"
)

func main() {
	// 连接服务器
	conn, err := net.Dial("tcp", "localhost:8080")
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

	// 发送消息
	message := TextMsg{
		Sender:   "me",
		Receiver: "you",
		Content:  "hello world",
	}
	msg, _ := json.Marshal(message)
	_, err = conn.Write(msg)
	if err != nil {
		Logger.Error("客户端发送消息失败", zap.Error(err))
		return
	}

	// 接收消息
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		Logger.Error("客户端接受消息异常", zap.Error(err))
		return
	}

	// 打印消息
	var msgReceive TextMsg
	err = json.Unmarshal(buf[:n], &msgReceive)
	color.Yellow("收到消息:%#v", msgReceive)
}
