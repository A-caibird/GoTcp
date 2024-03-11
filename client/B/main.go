package main

import (
	. "acaibird.com/zaplog"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"net"
)

func main() {
	// 连接服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
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
	_, err = conn.Write([]byte("hello world"))
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
	color.Yellow("收到消息:", string(buf[:n]))
}
