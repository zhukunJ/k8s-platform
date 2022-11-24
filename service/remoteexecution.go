package service

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

var Remoteexecution remoteexecution

type remoteexecution struct{}

var (
	client *ssh.Client
	err    error
)

//连接的配置
type ClientConfig struct {
	Host       string      //ip
	Port       int64       // 端口
	Username   string      //用户名
	Password   string      //密码
	Client     *ssh.Client //ssh client
	LastResult string      //最近一次运行的结果
}

//获取remoteexecution列表
func (p *remoteexecution) Remoteexecutions(app string) (stdout string, err error) {

	cliConf := new(ClientConfig)
	cliConf.createClient("114.55.233.102", 22, "root", "Wzzkj123")

	// fmt.Sprintf "容器当前状态" + cliConf.RunShell("docker ps -a")
	str := cliConf.RunShell("docker ps -a")

	return str, err
}

func (cliConf *ClientConfig) createClient(host string, port int64, username, password string) (err error) {

	cliConf.Host = host
	cliConf.Port = port
	cliConf.Username = username
	cliConf.Password = password
	cliConf.Port = port

	//一般传入四个参数：user，[]ssh.AuthMethod{ssh.Password(password)}, HostKeyCallback，超时时间，
	config := ssh.ClientConfig{
		User: cliConf.Username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)

	//获取client
	if client, err = ssh.Dial("tcp", addr, &config); err != nil {
		log.Fatalln("error occurred:", err)
	}

	cliConf.Client = client
	return nil
}

func (cliConf *ClientConfig) RunShell(shell string) string {
	var (
		session *ssh.Session
	)

	//获取session，这个session是用来远程执行操作的
	if session, err = cliConf.Client.NewSession(); err != nil {
		log.Fatalln("error occurred:", err)
	}

	//执行shell
	if output, err := session.CombinedOutput(shell); err != nil {
		log.Fatalln("error occurred:", err)
	} else {
		cliConf.LastResult = string(output)
	}
	return cliConf.LastResult
}
