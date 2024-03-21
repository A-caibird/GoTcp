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
			Logger.Error("断开与客户端链接异常:", zap.Error(err))
		}
		Logger.Info("客户端断开链接,goroutine退出!")
	}(conn)

	for {
		// 获取消息长度
		lens, err := getMsgLength(conn)
		if err != nil {
			return
		}

		// 获取消息字节数组
		byteMsg, err := getMsgBytesContent(conn, lens)

		// 消息解析
		var msgReceive message.TextMsg
		err = json.Unmarshal(byteMsg[:lens], &msgReceive)
		color.Blue("收到客户端消息:%#v", msgReceive)
		return
	}
}

func getMsgLength(conn net.Conn) (uint64, error) {
	buf := make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		if err == net.ErrClosed {
			Logger.Error("客户端关闭异常:", zap.Error(err))
			return 0, err
		} else if err == io.EOF {
			Logger.Error("客户端关闭异常:", zap.Error(err))
			return 0, err
		} else {
			Logger.Error("客户端读取消息长度错误:", zap.Error(err))
			return 0, err
		}
	}
	lens := binary.BigEndian.Uint64(buf[:n])
	return lens, nil
}

func getMsgBytesContent(conn net.Conn, lens uint64) ([]byte, error) {
	buf := make([]byte, lens)
	n, err := conn.Read(buf)
	if err != nil {
		Logger.Error("read error:", zap.Error(err))
	}
	return buf[:n], err
}
