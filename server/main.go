package main

import (
	. "acaibird.com/server/log"
	"acaibird.com/server/message"
	"acaibird.com/server/mysql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"io"
	"net"
)

var (
	MapUserConn = make(map[string]net.Conn)
)

func main() {
	// 监听端口
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		Logger.Error("Server listening port exception:", zap.Error(err))
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
			Logger.Error("Server accepts client connection error:", zap.Error(err))
			continue
		}

		// 处理连接
		go handleConn(conn)
	}
}

// 处理连接
func handleConn(conn net.Conn) {
	var client string
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			Logger.Error("An exception occurred while server trying to disconnect the client connection:", zap.Error(err), zap.String("client", client))
		}
		delete(MapUserConn, client)
		Logger.Info("Client disconnect with Server,goroutine exit!", zap.String("client", client))
	}(conn)

	for {
		// 获取消息长度
		lens, lenBytes, err := GetMsgLength(conn, client)
		if err != nil {
			return
		}

		// 获取消息字节数组
		byteMsg, err := GetMsgBytesContent(conn, lens, client)
		if err != nil {
			return
		}

		// 消息解析
		var msgReceive message.TextMsg
		err = json.Unmarshal(byteMsg[:lens], &msgReceive)
		if err != nil {
			Logger.Error("Client message parsing failed:", zap.Error(err), zap.String("client", client))
			continue
		}

		// 登录信息处理
		if msgReceive.Type == "login" {
			MapUserConn[msgReceive.Sender] = conn
			client = msgReceive.Sender
			fmt.Printf("User %s log in\n", msgReceive.Sender)

			// 发送离线消息
			msgs, _ := ReadTextMsgFromDB(msgReceive.Sender)
			for _, v := range msgs {
				color.Blue("%#v\n", v)
				if v.Receiver == msgReceive.Sender {
					SendOfflineTextMsg(conn, v)
					err := DelOfflineTextMsg(v.Receiver)
					err = mysqlDB.DBError(err)
					if err == nil {
						Logger.Info("Successfully pushed offline messages to the client, deleted offline messages successfully!", zap.Error(err), zap.String("client and receiver", client))
					}
				}
			}
			continue
		}
		// 普通文本消息
		if _, ok := MapUserConn[msgReceive.Receiver]; !ok {
			// 用户不在线,存储消息
			err = msgReceive.WriteToDB()
			err = mysqlDB.DBError(err)
			if err == nil {
				Logger.Info("Offline messages stored to the database successfully!", zap.Error(err), zap.String("from", msgReceive.Sender), zap.String("to", msgReceive.Receiver))
			}
		} else {
			//用户在线,发送消息
			_, err := MapUserConn[msgReceive.Receiver].Write(lenBytes)
			_, err = MapUserConn[msgReceive.Receiver].Write(byteMsg[:lens])
			if err != nil {
				Logger.Error("Message recipient is online, server push message exception:", zap.Error(err), zap.String("from", msgReceive.Sender), zap.String("to", msgReceive.Receiver))
			}
		}
	}
}

func SendOfflineTextMsg(conn net.Conn, msg message.TextMsg) {
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

func ReadTextMsgFromDB(receiver string) (msgs []message.TextMsg, err error) {
	db, err := mysqlDB.InitDB()
	if err != nil {
		return nil, errors.New("database connection exception")
	}
	defer func() {
		err = db.Close()
		if err != nil {
			return
		}
	}()

	// 准备 SQL 语句
	stmt, err := db.Prepare("SELECT type, sender, receiver, content, time FROM text_msgs WHERE receiver = ?")
	if err != nil {
		return nil, errors.New("pre-execute SQL statement exception")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	// 执行 SQL 语句
	rows, err := stmt.Query(receiver)
	if err != nil {
		return nil, errors.New("database executes SQL statement exception")
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			return
		}
	}()

	var msg message.TextMsg
	for rows.Next() {
		err = rows.Scan(&msg.Type, &msg.Sender, &msg.Receiver, &msg.Content, &msg.Time)
		if err == nil {
			msgs = append(msgs, msg)
		}
	}
	return msgs, nil
}

func DelOfflineTextMsg(receiver string) (err error) {
	db, err := mysqlDB.InitDB()
	if err != nil {
		return errors.New("database connection exception")
	}
	defer func() {
		err = db.Close()
		if err != nil {
			return // 如果数据库关闭失败,也返回了对应的err
		}

	}()
	stmt, err := db.Prepare("DELETE FROM text_msgs WHERE receiver = ?")
	if err != nil {
		return errors.New("pre-execute SQL statement exception")

	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()
	_, err = stmt.Exec(receiver)
	if err != nil {
		return errors.New("database executes SQL statement exception")
	}
	return err
}

func GetMsgLength(conn net.Conn, name string) (uint64, []byte, error) {
	buf := make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			Logger.Error("An exception occurred when client close ", zap.Error(err), zap.String("client", name))
			return 0, nil, err
		} else if err == io.EOF {
			Logger.Error("Read client message with incorrect length:", zap.Error(err), zap.String("client", name))
			return 0, nil, err
		}
	}
	// 读取到的字节数
	lens := binary.BigEndian.Uint64(buf[:n])
	return lens, buf, nil
}

func GetMsgBytesContent(conn net.Conn, lens uint64, name string) ([]byte, error) {
	buf := make([]byte, lens)
	// n 为读取到的字节数
	n, err := conn.Read(buf)
	if err != nil {
		Logger.Error("Read Message(type:[]byte) error:", zap.Error(err))
	}
	return buf[:n], err
}
