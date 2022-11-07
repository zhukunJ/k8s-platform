package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Node node

type node struct {}

//Node列表，支持过滤、排序、分页
func(p *node) GetNodes(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		FilterName string `form:"filter_name"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	data, err := service.Node.GetNodes(params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Node列表成功",
		"data": data,
	})
}
//Node詳情
func(p *node) GetNodeDetail(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		NodeName    string `form:"node_name"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	data, err := service.Node.GetNodeDetail(params.NodeName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	//測試
	//jsbyte, _ := json.Marshal(data)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Node詳情成功",
		"data": data,
	})
}