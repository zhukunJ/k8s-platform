package controller

import (
	"k8s-platform/mertics"

	"github.com/gin-gonic/gin"
)

//实例化router结构体，可使用该对象点出首字母大写的方法（跨包调用）
var Router router

//创建router结构体
type router struct{}

//初始化路由规则，创建测试api接口
func (r *router) InitApiRouter(router *gin.Engine) {
	router.
		//登录
		POST("/api/login", Login.Auth).
		//用户权限
		GET("/api/userInfo", Login.UserInfo).
		//登出
		GET("/api/logout", Login.Logout).
		//pod操作
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		DELETE("/api/k8s/pod/multiple", Pod.DeleteMultiplePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		GET("/api/k8s/pod/numns", Pod.GetPodNumPerNs).

		// mertics操作
		GET("/mertic", mertics.GetMertics).
		POST("/mertic/del", mertics.DeleteMertics).
		//deployment操作
		GET("/api/k8s/deployments", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		DELETE("/api/k8s/deployment/del", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).
		GET("/api/k8s/deployment/numns", Deployment.GetDeploymentNumPerNs).
		//daemonset操作
		GET("/api/k8s/daemonsets", DaemonSet.GetDaemonSets).
		GET("/api/k8s/daemonset/detail", DaemonSet.GetDaemonSetDetail).
		DELETE("/api/k8s/daemonset/del", DaemonSet.DeleteDaemonSet).
		PUT("/api/k8s/daemonset/update", DaemonSet.UpdateDaemonSet).
		//statefulset操作
		GET("/api/k8s/statefulsets", StatefulSet.GetStatefulSets).
		GET("/api/k8s/statefulset/detail", StatefulSet.GetStatefulSetDetail).
		DELETE("/api/k8s/statefulset/del", StatefulSet.DeleteStatefulSet).
		PUT("/api/k8s/statefulset/update", StatefulSet.UpdateStatefulSet).
		//service操作
		GET("/api/k8s/services", Servicev1.GetServicev1s).
		GET("/api/k8s/service/detail", Servicev1.GetServicev1Detail).
		DELETE("/api/k8s/service/del", Servicev1.DeleteServicev1).
		PUT("/api/k8s/service/update", Servicev1.UpdateServicev1).
		POST("/api/k8s/service/create", Servicev1.CreateServicev1).
		//POST("/api/k8s/service/create", Servicev1.CreateServicev1).
		//ingress操作
		GET("/api/k8s/ingresses", Ingress.GetIngresss).
		GET("/api/k8s/ingress/detail", Ingress.GetIngressDetail).
		DELETE("/api/k8s/ingress/del", Ingress.DeleteIngress).
		PUT("/api/k8s/ingress/update", Ingress.UpdateIngress).
		POST("/api/k8s/ingress/create", Ingress.CreateIngress).
		//configmap操作
		GET("/api/k8s/configmaps", ConfigMap.GetConfigMaps).
		GET("/api/k8s/configmap/detail", ConfigMap.GetConfigMapDetail).
		DELETE("/api/k8s/configmap/del", ConfigMap.DeleteConfigMap).
		PUT("/api/k8s/configmap/update", ConfigMap.UpdateConfigMap).
		//sercret操作
		GET("/api/k8s/secrets", Secret.GetSecrets).
		GET("/api/k8s/secret/detail", Secret.GetSecretDetail).
		DELETE("/api/k8s/secret/del", Secret.DeleteSecret).
		PUT("/api/k8s/secret/update", Secret.UpdateSecret).
		//pvc操作
		GET("/api/k8s/pvcs", Pvc.GetPvcs).
		GET("/api/k8s/pvc/detail", Pvc.GetPvcDetail).
		DELETE("/api/k8s/pvc/del", Pvc.DeletePvc).
		PUT("/api/k8s/pvc/update", Pvc.UpdatePvc).
		//node操作
		GET("/api/k8s/nodes", Node.GetNodes).
		GET("/api/k8s/node/detail", Node.GetNodeDetail).
		//namespace操作
		GET("/api/k8s/namespaces", Namespace.GetNamespaces).
		GET("/api/k8s/namespace/detail", Namespace.GetNamespaceDetail).
		DELETE("/api/k8s/namespace/del", Namespace.DeleteNamespace).
		//pv操作
		GET("/api/k8s/pvs", Pv.GetPvs).
		GET("/api/k8s/pv/detail", Pv.GetPvDetail)
	//workflow操作
	// GET("/api/k8s/workflows", Workflow.GetList).
	// GET("/api/k8s/workflow/detail", Workflow.GetById).
	// POST("/api/k8s/workflow/create", Workflow.Create).
	// DELETE("/api/k8s/workflow/del", Workflow.DelById)
}
