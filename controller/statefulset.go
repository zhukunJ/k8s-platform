package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var StatefulSet statefulSet

type statefulSet struct {}

//StatefulSet列表，支持过滤、排序、分页
func(p *statefulSet) GetStatefulSets(ctx *gin.Context) {
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
	data, err := service.StatefulSet.GetStatefulSets(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取StatefulSet列表成功",
		"data": data,
	})
}
//StatefulSet詳情
func(p *statefulSet) GetStatefulSetDetail(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		StatefulSetName    string `form:"statefulset_name"`
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
	data, err := service.StatefulSet.GetStatefulSetDetail(params.StatefulSetName, params.Namespace)
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
		"msg": "获取StatefulSet詳情成功",
		"data": data,
	})
}
//刪除StatefulSet
func(p *statefulSet) DeleteStatefulSet(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		StatefulSetName    string `json:"statefulset_name"`
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
	err := service.StatefulSet.DeleteStatefulSet(params.StatefulSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "刪除StatefulSet成功",
		"data": nil,
	})
}
//更新StatefulSet
func(p *statefulSet) UpdateStatefulSet(ctx *gin.Context) {
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
	err := service.StatefulSet.UpdateStatefulSet(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新StatefulSet成功",
		"data": nil,
	})
}