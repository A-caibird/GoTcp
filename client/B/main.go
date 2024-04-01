package main

import (
	"acaibird.com/clientA/message"
	. "acaibird.com/clientB/log"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"text/template"
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
			Logger.Error("客户端关闭异常!", zap.Error(err))
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
	// 消息长度
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(len(jsonLogin)))
	_, err = conn.Write(buf)       // 发送消息长度
	_, err = conn.Write(jsonLogin) // 发送消息内容

	for {
		ch := make(chan int, 1)
		go SendMessage(conn, &ch)
		go ReceiveMessage(conn)
		select {}
	}
}

func SendMessage(conn net.Conn, ch *chan int) {
	for {
		fmt.Printf("是否需要发送消息(y/n):  \n")
		var input string
		_, err2 := fmt.Scanln(&input)
		if err2 != nil {
			Logger.Error("输入异常", zap.Error(err2))
			return
		}
		if input == "n" {
			fmt.Println("退出消息发送!")
			*ch <- 1
			return // 退出
		}

		messageTypes := []message.MessageTypeTip{
			{1, "Text"},
			{2, "Image[png,jpg,jpeg,gif]"},
			{3, "Video[mp4]"},
		}
		tmpl := `请选择消息类型:
{{- range .}}
{{.Index}}:{{.Name}}
{{- end}}
`
		t := template.Must(template.New("messageTypes").Parse(tmpl))
		err := t.Execute(os.Stdout, messageTypes)
		if err != nil {
			return
		}

		// 输入消息类型
		reader := bufio.NewReader(os.Stdin)
		s, _ := reader.ReadString('\n')
		s = strings.TrimSpace(s)
		num, err := strconv.ParseInt(s, 10, 64)
		switch num {
		case 1:
			SendTextMsg(conn)
		case 2:
			fmt.Println("暂不支持图片消息!")
		case 3:
			fmt.Println("暂不支持视频消息!")
		default:
			fmt.Println("无效输入!")
		}
	}
}

func SendTextMsg(conn net.Conn) {
	msg := message.TextMsg{
		Type:     "text",
		Sender:   "A",
		Receiver: "B",
		Time:     time.Now(),
	}

	fmt.Printf("请输入接受者:\n")
	_, _ = fmt.Scanln(&msg.Receiver)
	fmt.Printf("请输入消息内容:\n")
	_, _ = fmt.Scanln(&msg.Content)

	buf := make([]byte, 8)
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		Logger.Error("A客户端发送消息,转化为json异常!", zap.Error(err))
		return
	}
	binary.BigEndian.PutUint64(buf, uint64(len(jsonMsg)))

	// 消息长度:字节数组
	_, err = conn.Write(buf)

	// 消息内容:字节数组
	_, err = conn.Write(jsonMsg)
	if err != nil {
		return
	}
}

func ReceiveMessage(conn net.Conn) {
	for {
		buf := make([]byte, 8)
		_, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				Logger.Info("disconnect server!", zap.Error(err))
				return
			} else if errors.Is(err, io.EOF) {
				Logger.Error("No more messages for now!", zap.Error(err))
				//continue
			}
		}
		lens := binary.BigEndian.Uint64(buf)
		var msg message.TextMsg
		buf = make([]byte, lens)
		_, err = conn.Read(buf)
		err = json.Unmarshal(buf, &msg)
		color.Blue("收到消息:%#v\n", msg)
	}
}
