package service

import (
	"github.com/wonderivan/logger"
	"k8s-platform/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

func(k *k8s) Init() {
	conf, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		logger.Error("创建k8s配置失败", err)
	}
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		logger.Error("创建k8s clientset失败", err)
	} else {
		logger.Info("创建k8s clientset成功")
	}
	//将初始化完成的clientset赋值给k8s结构体属性，用于全局调用
	k.ClientSet = clientset
}