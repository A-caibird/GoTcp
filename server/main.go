package main

import (
	. "acaibird.com/server/log"
	"acaibird.com/server/message"
	"encoding/binary"
	"encoding/json"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"io"
	"net"
)

func main() {
	// 监听端口
	listener, err := net.Listen("tcp", ":8081")
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

	for {
		// 获取消息长度
		buf := make([]byte, 8)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {

			Logger.Error("获取客户端消息字节长度错误:", zap.Error(err))
			return
		}
		lens := binary.BigEndian.Uint64(buf[:n])

		// 获取消息
		buf = make([]byte, lens)
		n, err = conn.Read(buf)
		if err != nil {
			Logger.Error("read error:", zap.Error(err))
		}

		var msgReceive message.TextMsg
		err = json.Unmarshal(buf[:n], &msgReceive)
		color.Blue("收到客户端消息:%#v", msgReceive)
	}
	// TODO: 当客户端关闭没有发送消息的时候应该怎么办
}
