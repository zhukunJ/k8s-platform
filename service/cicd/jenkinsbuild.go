package cicd

import (
	"log"
)

var JenkinsBuild jenkinsBuild

// 远程执行命令
type jenkinsBuild struct{}

func (p *jenkinsBuild) BuildJob(job string, params map[string]string) (string, error) {

	_, err := Jenkins.JenkinsClientSet.BuildJob(ctx, job, params)

	if err != nil {
		log.Println(err)
		return "", err
	}
	return "构建完成", nil
}

func (p *jenkinsBuild) GetResult(job string) (string, error) {

	app, err := Jenkins.JenkinsClientSet.GetJob(ctx, job)
	if err != nil {
		log.Println(err)
	}

	lastBuild, err := app.GetLastBuild(ctx)
	if err != nil {
		log.Println(err)
	}
	return lastBuild.GetResult(), nil
}
