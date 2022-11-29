package websocketflow

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"k8s-platform/service/cicd"

	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
)

var (
	ctx = context.Background()
)

// 根据传入参数build job
func GetJobDescribe(mjk *gojenkins.Jenkins, jobname string) (string, string, error) {
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
	params := new(struct {
		JobName string `form:"jobname"`
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

	go SendlogWeb(mySSH, params.JobName)

}

// 读取ssh输出，发送到websocket中
func SendlogWeb(mySSh *MySSH, jobname string) {
	// params := map[string]string{
	// 	"CHANGE_TYPE": "DEPLOY_PROD",
	// 	"GITBRACH":    "master",
	// }
	status, logs, err := GetJobDescribe(cicd.Jenkins.JenkinsClientSet, jobname)
	if err != nil {
		log.Printf("ERR, %v\n", err)
	}
	log.Println(status)
	jenlins_log_count := len(strings.Split(logs, "\n"))
	// if "End of Pipeline" in logs 就退出
	if strings.Contains(logs, "End of Pipeline") {
		for _, v := range strings.Split(logs, "\n") {
			// 读取ssh输出，发送到websocket中
			// 0.5秒
			time.Sleep(100 * time.Millisecond)
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
		// 基于for循环，每隔0.5秒，获取一次jenkins日志每次获取的日志长度
		for {
			//
			fmt.Println("for循环开始了!")
			time.Sleep(1000 * time.Millisecond)
			// 获取最后构建的日志
			status, logs, err := GetJobDescribe(cicd.Jenkins.JenkinsClientSet, jobname)
			if err != nil {
				log.Printf("ERR, %v\n", err)
			}
			log.Println(status)

			// 如果日志长度不一样，就发送基于jenkins_log_count的日志后长度，到websocket中
			if len(strings.Split(logs, "\n")) > jenlins_log_count {
				for _, v := range strings.Split(logs, "\n")[jenlins_log_count:len(strings.Split(logs, "\n"))] {
					// 读取ssh输出，发送到websocket中
					err := mySSh.Websocket.WriteMessage(websocket.TextMessage, []byte(v))
					if err != nil {
						fmt.Println("websocket发送数据失败：", err)
						break
					}
				}
				jenlins_log_count = len(strings.Split(logs, "\n"))
				fmt.Println("当前日志长度", jenlins_log_count)
				fmt.Println("我没发送完")
			}
			// 如果日志中包含"End of Pipeline"，就发送基于jenkins_log_count的日志后长度，到websocket中
			if strings.Contains(logs, "End of Pipeline") {
				fmt.Println("我即将结束啦")
				for _, v := range strings.Split(logs, "\n")[jenlins_log_count:] {
					// 读取ssh输出，发送到websocket中
					err := mySSh.Websocket.WriteMessage(websocket.TextMessage, []byte(v))
					if err != nil {
						fmt.Println("websocket发送数据失败：", err)
						mySSh.Websocket.Close()
						break
					}
				}
				fmt.Println("websocket关闭")
				break
			}
		}

	}
}
