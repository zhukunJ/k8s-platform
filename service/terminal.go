package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	"k8s-platform/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"time"
)

var Terminal terminal

type terminal struct {}

//websocket handler
func(t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	//加载k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		logger.Error("创建k8s配置失败，" + err.Error())
		return
	}
	//解析请求参数
	if err := r.ParseForm(); err != nil {
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	logger.Info("exec pod: %s, container: %s, namespace: %s\n", podName, containerName, namespace)
	//初始化terminalSession
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		logger.Error("get pty failed: %v\n", err)
		return
	}
	//处理关闭
	defer func() {
		logger.Info("close session.")
		pty.Close()
	}()
	req := K8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Container: containerName,
			Command:   []string{"/bin/bash"},
		}, scheme.ParameterCodec)
	logger.Info(req.URL())
	//升级SPDY协议
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		return
	}
	//定义SPDY的流传输
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err: %v \n", err)
		logger.Error(msg)
		pty.Write([]byte(msg))
		pty.Done()
	}
}

//定义终端和容器交互的内容格式
//Operation定义操作类型，stdin stdout
//Data具体的数据内容
//Rows和Cols可以理解为终端的行数和列数，也就是宽、高
type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

//初始化一个websocket.Upgrader实例对象，用于http协议升级为websocket协议
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

//实现ptyhandler，这个结构体中的websocket连接，接管传输内容
//wsConn是websocket的实例
//sizeChan用于定义终端的宽和高
//doneChan作为退出信号
type TerminalSession struct {
	wsConn *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

//New方法
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	//升级协议
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	//实例化TerminalSession对象
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return session, nil
}

//read方法，用于读取web端的输入
func(t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		log.Printf("read message err: %v", err)
		return 0, err
	}
	var msg TerminalMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		log.Printf("read parse mesage err: %v", err)
		return 0, err
	}

	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{
			Width:  msg.Cols,
			Height: msg.Rows,
		}
	case "ping":
		return 0, nil
	default:
		log.Printf("unknow message type %s", msg.Operation)
		return 0, fmt.Errorf("unknow message type %s", msg.Operation)
	}
	return 0, nil
}

//write方法，拿到容器执行的结果，写入web终端
func(t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("write parse message err: %v", err)
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("write message err: %v", err)
		return 0, err
	}
	return len(p), nil
}

//done,发出关闭信号
func(t *TerminalSession) Done() {
	close(t.doneChan)
}
//next,获取终端是否resize，或者是否退出
func(t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <- t.sizeChan:
		return &size
	case <- t.doneChan:
		return nil
	}
}
//close
func(t *TerminalSession) Close()  error {
	return t.wsConn.Close()
}