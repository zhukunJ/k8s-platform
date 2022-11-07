package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Pvc pvc

type pvc struct {}

//Pvc列表，支持过滤、排序、分页
func(p *pvc) GetPvcs(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
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
	data, err := service.Pvc.GetPvcs(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Pvc列表成功",
		"data": data,
	})
}
//Pvc詳情
func(p *pvc) GetPvcDetail(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		PvcName    string `form:"pvc_name"`
		Namespace  string `form:"namespace"`
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
	data, err := service.Pvc.GetPvcDetail(params.PvcName, params.Namespace)
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
		"msg": "获取Pvc詳情成功",
		"data": data,
	})
}
//刪除Pvc
func(p *pvc) DeletePvc(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		PvcName    string `json:"pvc_name"`
		Namespace  string `json:"namespace"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	err := service.Pvc.DeletePvc(params.PvcName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "刪除Pvc成功",
		"data": nil,
	})
}
//更新Pvc
func(p *pvc) UpdatePvc(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		Content    string `json:"content"`
		Namespace  string `json:"namespace"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}

	//调用service方法，获取数据
	err := service.Pvc.UpdatePvc(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新Pvc成功",
		"data": nil,
	})
}