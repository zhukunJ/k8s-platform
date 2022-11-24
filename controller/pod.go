package controller

import (
	"encoding/json"
	"io/ioutil"
	"k8s-platform/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/wonderivan/logger"
)

var Pod pod

type pod struct{}

//Pod列表，支持过滤、排序、分页
func (p *pod) GetPods(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"` // 如果前端请求没有namespace参数默认返回所有命名空间的pods
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod列表成功",
		"data": data,
	})
}

//Pod詳情
func (p *pod) GetPodDetail(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod詳情成功",
		"data": data,
	})
}

//刪除Pod
func (p *pod) DeletePod(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "刪除Pod成功",
		"data": nil,
	})
}

//刪除多个Pod
func (p *pod) DeleteMultiplePod(ctx *gin.Context) {
	pods := ctx.Request.Body
	body, _ := ioutil.ReadAll(pods)
	data := gjson.Get(string(body), "data")
	// data返回数据为[{"podname":"nginx-96b7dc68-fxjwl","namespace":"default"},{"podname":"storage-provisioner","namespace":"kube-system"}]
	var podList service.PodList
	err := json.Unmarshal([]byte(data.String()), &podList)
	if err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	// fmt.Println(podList)
	//调用service方法，获取数据
	err = service.Pod.DeleteMultiplePod(podList)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "批量刪除成功",
		"data": nil,
	})
}

//更新Pod
func (p *pod) UpdatePod(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	err := service.Pod.UpdatePod(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新Pod成功",
		"data": nil,
	})
}

//獲取Pod的容器名
func (p *pod) GetPodContainer(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	data, err := service.Pod.GetPodContainer(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "獲取Pod容器名成功",
		"data": data,
	})
}

//獲取Pod日誌
func (p *pod) GetPodLog(ctx *gin.Context) {
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		ContainerName string `form:"container_name"`
		PodName       string `form:"pod_name"`
		Namespace     string `form:"namespace"`
	})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败,", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取数据
	data, err := service.Pod.GetPodLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "獲取Pod日誌成功",
		"data": data,
	})
}

//獲取每個命名空間Pod數量
func (p *pod) GetPodNumPerNs(ctx *gin.Context) {
	//调用service方法，获取数据
	data, err := service.Pod.GetPodNumPerNs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "獲取Pod數量成功",
		"data": data,
	})
}
