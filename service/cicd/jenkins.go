package cicd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var JenkinsBuild jenkinsBuild

// 远程执行命令
type jenkinsBuild struct{}

// 将job.GetName()和lastBuild.GetResult()放到结构体中

type DataTwo struct {
	JobName string `json:"jobName"`
	Status  string `json:"result"`
}

type DataOne struct {
	Data []DataTwo `json:"data"`
}
type DeploymentsResp struct {
	Items []DataTwo `json:"items"`
}

type Job struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
type Js struct {
	Jobs []Job `json:"jobs"`
}

func (p *jenkinsBuild) BuildJob(job string, params map[string]string) (string, error) {

	_, err := Jenkins.JenkinsClientSet.BuildJob(ctx, job, params)

	if err != nil {
		log.Println(err)
		return "", err
	}
	return "构建完成", nil
}

func (p *jenkinsBuild) GetJobAll() (DeploymentsResp, int, error) {
	var NewStatus DataOne

	url := "http://admin:wzzkj123@101.42.13.214:8080/api/json?pretty=true"
	// 请求url 获取数据
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// 解析数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//3.反序列化
	var jenkins Js
	err2 := json.Unmarshal([]byte(body), &jenkins)
	if err2 != nil {
		fmt.Printf("unmarshal failed, err:%v", err)
	}
	for _, job := range jenkins.Jobs {
		if job.Color == "blue" {
			NewStatus.Data = append(NewStatus.Data, DataTwo{JobName: job.Name, Status: "成功"})
		} else if job.Color == "red" {
			NewStatus.Data = append(NewStatus.Data, DataTwo{JobName: job.Name, Status: "失败"})
		} else if job.Color == "blue_anime" {
			NewStatus.Data = append(NewStatus.Data, DataTwo{JobName: job.Name, Status: "构建中"})
		} else {
			NewStatus.Data = append(NewStatus.Data, DataTwo{JobName: job.Name, Status: "未知"})
		}
	}

	da := DeploymentsResp{
		Items: NewStatus.Data,
	}
	total := len(da.Items)

	return da, total, nil
}
