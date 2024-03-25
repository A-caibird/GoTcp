package main

import (
	. "acaibird.com/server/log"
	"acaibird.com/server/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
)

var (
	MapUserMsg  = make(map[string]message.TextMsg)
	MapUserConn = make(map[string]net.Conn)
)

func main() {
	// 监听端口
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		Logger.Error("服务器监听端口异常:", zap.Error(err))
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
		}
	}(listener)

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
		for v, i := range MapUserMsg {
			fmt.Printf("用户:%s,消息:%#v\n", v, i)
		}

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
		if err != nil {
			Logger.Error("消息解析异常:", zap.Error(err))
			return
		}
		fmt.Printf("%#v\n", msgReceive)

		// 登录信息处理
		if msgReceive.Type == "login" {
			MapUserConn[msgReceive.Sender] = conn
			fmt.Printf("用户%s上线\n", msgReceive.Sender)

			// 发送离线消息
			for v, _ := range MapUserMsg {
				if v == msgReceive.Sender {

					fmt.Println("推送离线消息给用户:", v)
					sendOfflineMsg(conn, MapUserMsg[v])
					delete(MapUserMsg, v)
				}
			}
			continue
		}

		// 普通文本消息
		if _, ok := MapUserConn[msgReceive.Receiver]; !ok {
			// 用户不在线,存储消息
			MapUserMsg[msgReceive.Receiver] = msgReceive
			fmt.Printf("用户%s不在线,存储消息:  \n", msgReceive.Receiver)
		} else {
			// 用户在线,发送消息
			_, err2 := MapUserConn[msgReceive.Receiver].Write(byteMsg[:lens])
			if err2 != nil {
				return
			}
		}
		return
	}
}

func sendOfflineMsg(conn net.Conn, msg message.TextMsg) {
	buf := make([]byte, 8)
	jsonMsg, _ := json.Marshal(msg)
	binary.BigEndian.PutUint64(buf, uint64(len(jsonMsg)))

	// 消息长度
	_, err := conn.Write(buf)

	// 消息内容
	_, err = conn.Write(jsonMsg)
	if err != nil {
		return
	}
}

func getMsgLength(conn net.Conn) (uint64, error) {
	buf := make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
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
	//fmt.Println("接收到消息长度:", n)
	return buf[:n], err
}
