package main

import (
	"acaibird.com/clientA/message"
	. "acaibird.com/clientB/log"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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

	var wg sync.WaitGroup

	// A发送消息
	wg.Add(1)
	go SendMessage(conn, &wg)

	// A接受消息
	wg.Add(1)
	go ReceiveMessage(conn, &wg)

	// 退出
	wg.Wait()
	os.Exit(0)
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
	jsonMsg, _ := json.Marshal(msg)
	binary.BigEndian.PutUint64(buf, uint64(len(jsonMsg)))

	// 消息长度:字节数组
	_, err := conn.Write(buf)

	// 消息内容:字节数组
	_, err = conn.Write(jsonMsg)
	if err != nil {
		return
	}
}

func SendMessage(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Printf("是否需要发送消息(y/n):  ")
		var input string
		_, err2 := fmt.Scanln(&input)
		if err2 != nil {
			Logger.Error("输入异常", zap.Error(err2))
			return
		}
		if input == "n" {
			fmt.Println("退出消息发送!")
			return // 退出
		}

		messageTypes := []message.MessageType{
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
			fmt.Println("暂不支持图片消息")
		case 3:
			fmt.Println("暂不支持视频消息")
		default:
			fmt.Println("输入有误")
		}
	}
}

func ReceiveMessage(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		buf := make([]byte, 8)
		_, err := conn.Read(buf)
		if err != nil {
			return
		}
		lens := binary.BigEndian.Uint64(buf)
		var msg message.TextMsg
		buf = make([]byte, lens)
		_, err = conn.Read(buf)
		err = json.Unmarshal(buf, &msg)
		color.Blue("收到消息:%#v\n", msg)
	}
}
