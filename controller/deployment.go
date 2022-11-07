package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Deployment deployment

type deployment struct {}

//Deployment列表，支持过滤、排序、分页
func(p *deployment) GetDeployments(ctx *gin.Context) {
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
	data, err := service.Deployment.GetDeployments(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Deployment列表成功",
		"data": data,
	})
}
//Deployment詳情
func(p *deployment) GetDeploymentDetail(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		DeploymentName    string `form:"deployment_name"`
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
	data, err := service.Deployment.GetDeploymentDetail(params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	//測試
	jsbyte, _ := json.Marshal(data)
	a := string(jsbyte)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Deployment詳情成功",
		"data": a,
	})
	//ctx.JSON(http.StatusOK, gin.H{
	//	"msg": "获取Deployment詳情成功",
	//	"data": data,
	//})
}
//刪除Deployment
func(p *deployment) DeleteDeployment(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		DeploymentName    string `json:"deployment_name"`
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
	err := service.Deployment.DeleteDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "刪除Deployment成功",
		"data": nil,
	})
}
//更新Deployment
func(p *deployment) UpdateDeployment(ctx *gin.Context) {
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
	err := service.Deployment.UpdateDeployment(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新Deployment成功",
		"data": nil,
	})
}
//修改Deployment副本數
func(p *deployment) ScaleDeployment(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		ScaleNum          int    `json:"scale_num"`
		DeploymentName    string `json:"deployment_name"`
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
	replicas, err := service.Deployment.ScaleDeployment(params.DeploymentName, params.Namespace, params.ScaleNum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": replicas,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改Deployment副本數成功",
		"data": nil,
	})
}
//重啟Deployment
func(p *deployment) RestartDeployment(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct{
		DeploymentName    string `json:"deployment_name"`
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
	err := service.Deployment.RestartDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "重啟Deployment成功",
		"data": nil,
	})
}
//創建Deployment
func(p *deployment) CreateDeployment(ctx *gin.Context) {
	deployCreate := &service.DeployCreate{}
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.ShouldBindJSON(deployCreate); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	err := service.Deployment.CreateDeployment(deployCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "創建Deployment成功",
		"data": nil,
	})
}
//獲取每個命名空間Deployment數量
func(p *deployment) GetDeploymentNumPerNs(ctx *gin.Context) {
	//调用service方法，获取数据
	data, err := service.Deployment.GetDeploymentNumPerNs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "獲取Deployment數量成功",
		"data": data,
	})
}