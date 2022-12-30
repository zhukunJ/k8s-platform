package cicd

import (
	"context"
	"log"

	"github.com/bndr/gojenkins"
)

var Jenkins jenkins

var (
	ctx = context.Background()
)

type jenkins struct {
	JenkinsClientSet *gojenkins.Jenkins
	Url              string
	User             string
	Password         interface{}
}

func (j *jenkins) Init() {
	client := &jenkins{
		Url:      "http://101.42.13.214:8080/",
		User:     "admin",
		Password: "admin",
	}

	jenkins := gojenkins.CreateJenkins(nil, client.Url, client.User, client.Password)
	_, err := jenkins.Init(ctx)

	if err != nil {
		log.Printf("jenkins初始化失败, %v\n", err)
	}
	log.Println("Jenkins UP")
	j.JenkinsClientSet = jenkins
}
