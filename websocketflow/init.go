package websocketflow

import (
	"bytes"
	"io"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// 定义一个结构体 方便保存各种连接信息
type MySSH struct {
	Websocket *websocket.Conn
	Stdin     io.WriteCloser
	Stdout    *wsBufferWriter
	Session   *ssh.Session
}

// 定义一个wsBufferWriter 并且写入时候加锁 防止stdout跟stderr同时写入
type wsBufferWriter struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}
