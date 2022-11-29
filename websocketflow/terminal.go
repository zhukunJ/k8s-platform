package websocketflow

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin" // go get -u github.com/gin-gonic/gin
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	"golang.org/x/crypto/ssh"
)

//定义write方法， 防止stdout跟stderr同时写入
func (w *wsBufferWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

// 程序入口
func RunWebSSH(c *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		Ip string `form:"ip"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := c.Bind(params); err != nil {
		logger.Error("参数绑定失败,", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	mySSH := &MySSH{}

	// 1. 升级请求websocket
	upGrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024 * 1024 * 10,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"webssh"},
	}

	webcon, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("升级http 为websoket失败：", err)
	}
	mySSH.Websocket = webcon // 将websocket连接保存到对象中

	// 创建一个ssh的配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 100, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		Auth: []ssh.AuthMethod{ssh.Password("Wzzkj123")},
	}
	ipddress := params.Ip + ":22"
	sshClient, err := ssh.Dial("tcp", ipddress, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	session, err := sshClient.NewSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	mySSH.Session = session

	// 保存输入流
	mySSH.Stdin, err = session.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	//保存ssh输出流
	sshOut := new(wsBufferWriter)
	session.Stdout = sshOut
	session.Stderr = sshOut
	mySSH.Stdout = sshOut

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 30, 120, modes); err != nil {
		fmt.Println("绑定pty失败：", err)
		return
	}

	session.Shell()

	//执行远程命令
	go Send2SSH(mySSH)
	go Send2Web(mySSH)

}

// 读取websocket数据，发送到ssh输入流中
func Send2SSH(mySSh *MySSH) {
	for {
		//read websocket msg  需要通过msgType 判断是传输类型
		_, wsData, err := mySSh.Websocket.ReadMessage()
		if err != nil {
			fmt.Println("读取websocket数据失败：", err)
			return
		}
		_, err = mySSh.Stdin.Write(wsData)
		if err != nil {
			fmt.Println("ssh发送数据失败：", err)
		}
		// fmt.Println("ssh发送数据：", string(wsData))

	}

}

// 读取ssh输出，发送到websocket中
func Send2Web(mySSh *MySSH) {
	for {
		if mySSh.Stdout.buffer.Len() > 0 {
			err := mySSh.Websocket.WriteMessage(websocket.TextMessage, mySSh.Stdout.buffer.Bytes())
			if err != nil {
				fmt.Println("websocket发送数据失败：", err)
			}
			mySSh.Stdout.buffer.Reset() //读完清空
		}

	}
}
