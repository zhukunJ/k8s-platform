package controller

import (
	"net/http"

	"k8s-platform/service/cicd"

	"github.com/gin-gonic/gin"
)

var JenkinsBuild jenkinsbuild

// 远程执行命令
type jenkinsbuild struct{}

func (p *jenkinsbuild) BuildJob(ctx *gin.Context) {
	// params := new(struct {
	// 	Name string `form:"name"`
	// })
	// if err := ctx.Bind(params); err != nil {
	// 	logger.Error("参数绑定失败,", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"msg":  err.Error(),
	// 		"data": nil,
	// 	})
	// 	return
	// }
	params := map[string]string{
		"CHANGE_TYPE": "DEPLOY_PROD",
		"GITBRACH":    "master",
	}
	data, err := cicd.JenkinsBuild.BuildJob("demo", params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "install success",
		"data": data,
	})
}

func (p *jenkinsbuild) GetResult(ctx *gin.Context) {
	// params := new(struct {
	// 	Name string `form:"name"`
	// })
	// if err := ctx.Bind(params); err != nil {
	// 	logger.Error("参数绑定失败,", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"msg":  err.Error(),
	// 		"data": nil,
	// 	})
	// 	return
	// }

	result, err := cicd.JenkinsBuild.GetResult("demo")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "上一次构建结果",
		"data": result,
	})
}
