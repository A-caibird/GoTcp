package main

import (
	"fmt"
	"net"
)

func main() {
	// 监听端口
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	defer listener.Close()

	// 循环接受连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
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
			fmt.Println("read error:", err)
			return
		}

		// 打印消息
		fmt.Println("收到消息:", string(buf[:n]))

		// 发送消息
		_, err = conn.Write([]byte("hello world"))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}
