package main

import (
	. "acaibird.com/zaplog"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"net"
)

func main() {
	// 监听端口
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		Logger.Error("服务器监听端口异常:", zap.Error(err))
		return
	}
	defer listener.Close()

	// 循环接受连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			Logger.Error("accept error:", zap.Error(err))
			continue
		}

		// 处理连接
		go handleConn(conn)
	}
}

// 处理连接
func handleConn(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// 接收消息
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			Logger.Error("read error:", zap.Error(err))
			return
		}

		// 打印消息
		color.Blue("收到客户端消息:", string(buf[:n]))

		// 发送消息
		_, err = conn.Write([]byte("hello world"))
		if err != nil {
			Logger.Error("推送消息失败:", zap.Error(err))
			return
		}
	}
}
