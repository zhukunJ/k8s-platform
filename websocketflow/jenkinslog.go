package websocketflow

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	ctx = context.Background()
)

type Intelops struct {
	Url      string
	User     string
	Password interface{}
}

// 根据传入参数build job
func (i *Intelops) BuildJob(mjk *gojenkins.Jenkins, jobname string, params map[string]string) (string, string, error) {
	// _, err := mjk.BuildJob(ctx, jobname, params)
	// if err != nil {
	// 	log.Println(err)
	// }
	//获取job的状态

	job, err := mjk.GetJob(ctx, jobname)
	if err != nil {
		log.Println(err)
	}
	lastBuild, err := job.GetLastBuild(ctx)
	if err != nil {
		log.Println(err)
	}

	// 获取最后构建的日志
	logs := lastBuild.GetConsoleOutput(ctx)

	//lastBuild.GetResult() :获取最后一次构建的状态
	return lastBuild.GetResult(), logs, err
}

func LogIndex(logs string) string {
	return logs[strings.Index(logs, "构建线上镜像 ok!")+len("构建线上镜像 ok!"):]

}

// 程序入口
func RunWebLog(c *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	// params := new(struct {
	// 	Ip string `form:"ip"`
	// })
	// //绑定参数，给匿名结构体中的属性赋值，值是入参
	// //form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	// if err := c.Bind(params); err != nil {
	// 	logger.Error("参数绑定失败,", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"msg":  err.Error(),
	// 		"data": nil,
	// 	})
	// 	return
	// }

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

	go SendlogWeb(mySSH)

}

// 读取ssh输出，发送到websocket中
func SendlogWeb(mySSh *MySSH) {
	client := &Intelops{
		Url:      "http://114.55.233.102:8080/",
		User:     "admin",
		Password: "admin",
	}

	jenkins := gojenkins.CreateJenkins(nil, client.Url, client.User, client.Password)
	_, err := jenkins.Init(ctx)

	if err != nil {
		log.Printf("ERR, %v\n", err)
	}
	log.Println("Jenkins UP")

	// build demo job and params : CHANGE_TYPE = DEPLOY_PROD ,GITBRACH = master 基于匿名结构体
	params := map[string]string{
		"CHANGE_TYPE": "DEPLOY_PROD",
		"GITBRACH":    "master",
	}
	status, logs, err := client.BuildJob(jenkins, "demo", params)
	if err != nil {
		log.Printf("ERR, %v\n", err)
	}
	log.Println(status)
	jenlins_log_count := len(strings.Split(logs, "\n"))
	// if "End of Pipeline" in logs 就退出
	if strings.Contains(logs, "End of Pipeline") {
		for _, v := range strings.Split(logs, "\n") {
			// 读取ssh输出，发送到websocket中
			err := mySSh.Websocket.WriteMessage(websocket.TextMessage, []byte(v))
			if err != nil {
				fmt.Println("websocket发送数据失败：", err)
				mySSh.Websocket.Close()

				break
			}
		}
		fmt.Println("websocket关闭")
	} else {
		// 先发送一次
		for _, v := range strings.Split(logs, "\n") {
			// 读取ssh输出，发送到websocket中
			err := mySSh.Websocket.WriteMessage(websocket.TextMessage, []byte(v))
			if err != nil {
				fmt.Println("websocket发送数据失败：", err)
				// mySSh.Websocket.Close()
				// break
			}
		}
		// 匿名函数jenkins_log_repeat

		str := jenkins_log_repeat(jenlins_log_count, mySSh, client, jenkins, params)
		fmt.Println(str)

	}
}

func jenkins_log_repeat(jenlins_log_count int, mySSh *MySSH, client *Intelops, jenkins *gojenkins.Jenkins, params map[string]string) string {
	_, jenlins_log_new, err := client.BuildJob(jenkins, "demo", params)
	if err != nil {
		log.Printf("ERR, %v\n", err)
	}
	jenkins_new_count := len(strings.Split(jenlins_log_new, "\n"))
	if strings.Contains(jenlins_log_new, "End of Pipeline") {
		for _, v := range strings.Split(jenlins_log_new, "\n")[jenlins_log_count:] {
			// 读取ssh输出，发送到websocket中
			err := mySSh.Websocket.WriteMessage(websocket.TextMessage, []byte(v))
			if err != nil {
				fmt.Println("websocket发送数据失败：", err)
				mySSh.Websocket.Close()
				return "递归结束"

			}
		}

	} else {
		jenkins_log_list_new := strings.Split(jenlins_log_new, "\n")[jenlins_log_count:jenkins_new_count]
		// 如果jenkins_log_list_new为空
		if len(jenkins_log_list_new) == 0 {
			fmt.Println("jenkins_log_list_new为空")

		} else {
			for _, v := range jenkins_log_list_new {
				// 读取ssh输出，发送到websocket中
				err := mySSh.Websocket.WriteMessage(websocket.TextMessage, []byte(v))
				if err != nil {
					fmt.Println("websocket发送数据失败：", err)
					mySSh.Websocket.Close()
					return "递归结束"
				}
			}
		}

	}
	jenkins_log_repeat(jenkins_new_count, mySSh, client, jenkins, params)
	return "递归结束"

}
