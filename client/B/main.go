package main

import (
	"fmt"
	"net"
)

func main() {
	// 连接服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	// 发送消息
	_, err = conn.Write([]byte("hello world"))
	if err != nil {
		fmt.Println("write error:", err)
		return
	}

	// 接收消息
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}

	// 打印消息
	fmt.Println("收到消息:", string(buf[:n]))
}
