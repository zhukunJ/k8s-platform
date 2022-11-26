package router

import (
	"k8s-platform/mertics"
	"k8s-platform/websocketflow"

	"k8s-platform/controller"

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
		POST("/api/login", controller.Login.Auth).
		//用户权限
		GET("/api/userInfo", controller.Login.UserInfo).
		//登出
		GET("/api/logout", controller.Login.Logout).
		//pod操作
		GET("/api/k8s/pods", controller.Pod.GetPods).
		GET("/api/k8s/pod/detail", controller.Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", controller.Pod.DeletePod).
		DELETE("/api/k8s/pod/multiple", controller.Pod.DeleteMultiplePod).
		PUT("/api/k8s/pod/update", controller.Pod.UpdatePod).
		GET("/api/k8s/pod/container", controller.Pod.GetPodContainer).
		GET("/api/k8s/pod/log", controller.Pod.GetPodLog).
		GET("/api/k8s/pod/numns", controller.Pod.GetPodNumPerNs).

		// mertics操作
		GET("/mertic", mertics.GetMertics).
		POST("/mertic/del", mertics.DeleteMertics).

		// 远程执行命令
		POST("/api/host/remoteexecution", controller.Remoteexecution.GetRemoteexecutions).

		//websocket host
		GET("/api/host/ws", websocketflow.RunWebSSH).

		//websocket jenkinslogs
		GET("/api/jenkinslogs/ws", websocketflow.RunWebLog).

		//deployment操作
		GET("/api/k8s/deployments", controller.Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", controller.Deployment.GetDeploymentDetail).
		DELETE("/api/k8s/deployment/del", controller.Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment/update", controller.Deployment.UpdateDeployment).
		PUT("/api/k8s/deployment/restart", controller.Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/scale", controller.Deployment.ScaleDeployment).
		POST("/api/k8s/deployment/create", controller.Deployment.CreateDeployment).
		GET("/api/k8s/deployment/numns", controller.Deployment.GetDeploymentNumPerNs).
		//daemonset操作
		GET("/api/k8s/daemonsets", controller.DaemonSet.GetDaemonSets).
		GET("/api/k8s/daemonset/detail", controller.DaemonSet.GetDaemonSetDetail).
		DELETE("/api/k8s/daemonset/del", controller.DaemonSet.DeleteDaemonSet).
		PUT("/api/k8s/daemonset/update", controller.DaemonSet.UpdateDaemonSet).
		//statefulset操作
		GET("/api/k8s/statefulsets", controller.StatefulSet.GetStatefulSets).
		GET("/api/k8s/statefulset/detail", controller.StatefulSet.GetStatefulSetDetail).
		DELETE("/api/k8s/statefulset/del", controller.StatefulSet.DeleteStatefulSet).
		PUT("/api/k8s/statefulset/update", controller.StatefulSet.UpdateStatefulSet).
		//service操作
		GET("/api/k8s/services", controller.Servicev1.GetServicev1s).
		GET("/api/k8s/service/detail", controller.Servicev1.GetServicev1Detail).
		DELETE("/api/k8s/service/del", controller.Servicev1.DeleteServicev1).
		PUT("/api/k8s/service/update", controller.Servicev1.UpdateServicev1).
		POST("/api/k8s/service/create", controller.Servicev1.CreateServicev1).
		//POST("/api/k8s/service/create", Servicev1.CreateServicev1).
		//ingress操作
		GET("/api/k8s/ingresses", controller.Ingress.GetIngresss).
		GET("/api/k8s/ingress/detail", controller.Ingress.GetIngressDetail).
		DELETE("/api/k8s/ingress/del", controller.Ingress.DeleteIngress).
		PUT("/api/k8s/ingress/update", controller.Ingress.UpdateIngress).
		POST("/api/k8s/ingress/create", controller.Ingress.CreateIngress).
		//configmap操作
		GET("/api/k8s/configmaps", controller.ConfigMap.GetConfigMaps).
		GET("/api/k8s/configmap/detail", controller.ConfigMap.GetConfigMapDetail).
		DELETE("/api/k8s/configmap/del", controller.ConfigMap.DeleteConfigMap).
		PUT("/api/k8s/configmap/update", controller.ConfigMap.UpdateConfigMap).
		//sercret操作
		GET("/api/k8s/secrets", controller.Secret.GetSecrets).
		GET("/api/k8s/secret/detail", controller.Secret.GetSecretDetail).
		DELETE("/api/k8s/secret/del", controller.Secret.DeleteSecret).
		PUT("/api/k8s/secret/update", controller.Secret.UpdateSecret).
		//pvc操作
		GET("/api/k8s/pvcs", controller.Pvc.GetPvcs).
		GET("/api/k8s/pvc/detail", controller.Pvc.GetPvcDetail).
		DELETE("/api/k8s/pvc/del", controller.Pvc.DeletePvc).
		PUT("/api/k8s/pvc/update", controller.Pvc.UpdatePvc).
		//node操作
		GET("/api/k8s/nodes", controller.Node.GetNodes).
		GET("/api/k8s/node/detail", controller.Node.GetNodeDetail).
		//namespace操作
		GET("/api/k8s/namespaces", controller.Namespace.GetNamespaces).
		GET("/api/k8s/namespace/detail", controller.Namespace.GetNamespaceDetail).
		DELETE("/api/k8s/namespace/del", controller.Namespace.DeleteNamespace).
		//pv操作
		GET("/api/k8s/pvs", controller.Pv.GetPvs).
		GET("/api/k8s/pv/detail", controller.Pv.GetPvDetail)
	//workflow操作
	// GET("/api/k8s/workflows", controller.Workflow.GetList).
	// GET("/api/k8s/workflow/detail", controller.Workflow.GetById).
	// POST("/api/k8s/workflow/create", controller.Workflow.Create).
	// DELETE("/api/k8s/workflow/del", controller.Workflow.DelById)
}
