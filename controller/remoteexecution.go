package controller

import (
	"net/http"

	"k8s-platform/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Remoteexecution remoteexecution

// 远程执行命令
type remoteexecution struct{}

func (p *remoteexecution) GetRemoteexecutions(ctx *gin.Context) {
	params := new(struct {
		Name string `form:"name"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Remoteexecution.Remoteexecutions(params.Name)
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
