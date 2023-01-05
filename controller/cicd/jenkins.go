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
	params := new(struct {
		Name       string `form:"name"`
		Changetype string `form:"changetype"`
		Branch     string `form:"branch"`
	})
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	params2 := map[string]string{
		"CHANGE_TYPE": params.Changetype,
		"GITBRACH":    params.Branch,
	}
	data, err := cicd.JenkinsBuild.BuildJob(params.Name, params2)
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

func (p *jenkinsbuild) GetJobAll(ctx *gin.Context) {
	result, total, err := cicd.JenkinsBuild.GetJobAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "上一次构建结果",
		"data":  result,
		"total": total,
	})
}
