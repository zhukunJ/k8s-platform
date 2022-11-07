package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct {}

//定义deployment列表的返回内容，Items是deployment列表，Total为deployment元素总数
//先过滤，再拿total，再做分页
type DeploymentsResp struct{
	Items []appsv1.Deployment `json:"items"`
	Total int         `json:"total"`
}

type DeploymentsNs struct {
	Namespace string `json:"namespace"`
	DeploymentNum    int    `json:"deployment_num"`
}

//定義結構體，用於創建deployment
type DeployCreate struct {
	Name          string   `json:"name"`
	Namespace     string   `json:"namespace"`
	Replicas      int32    `json:"replicas"`
	Image         string   `json:"image"`
	Label         map[string]string  `json:"label"`
	Cpu           string   `json:"cpu"`
	Memory        string   `json:"memory"`
	ContainerPort int32    `json:"container_port"`
	HealthCheck   bool     `json:"health_check"`
	HealthPath    string   `json:"health_path"`
}

//获取deployment列表
func(p *deployment) GetDeployments(filterName, namespace string, limit, page int) (deploymentsResp *DeploymentsResp, err error) {
	//通过clintset获取deployment完整列表
	deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取deployment列表失败", err)
		return nil, errors.New("获取deployment列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(deploymentList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery:   &Filter{Name:filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤
	filtered := selectableData.Filter()
	//再拿Total
	total := len(filtered.GenericDataList)
	//再排序和分页
	data := filtered.Sort().Paginate()
	//再将DataCell切片数据转成原生deployment切片
	deployments := p.fromCells(data.GenericDataList)
	//返回
	return &DeploymentsResp{
		Items: deployments,
		Total: total,
	}, nil
}
//获取deployment详情
func(p *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Deployment详情失败", err)
		return nil, errors.New("获取Deployment详情失败" + err.Error())
	}
	return deployment, nil
}
//删除deployment
func(p *deployment) DeleteDeployment(deploymentName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Deployment失败", err)
		return errors.New("删除Deployment失败" + err.Error())
	}
	return nil
}
//更新deployment
func(d *deployment) UpdateDeployment(namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}

	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Deployment失败, " + err.Error()))
		return errors.New("更新Deployment失败, " + err.Error())
	}
	return nil
}
//修改deployment副本數
func(p *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replicas int32, err error) {
	//獲取autoscaling.Scale類型的對象，能點出當前的副本數
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("獲取Deployment副本數失败", err)
		return 0, errors.New("獲取Deployment副本數失败" + err.Error())
	}
	//修改副本數
	scale.Spec.Replicas = int32(scaleNum)
	//更新副本數
	newScale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Deployment副本數失败", err)
		return 0, errors.New("更新Deployment副本數失败" + err.Error())
	}
	return newScale.Spec.Replicas, nil
}
//重啟Deployment
func(p *deployment) RestartDeployment(deploymentName, namespace string) (err error) {
	//使用patchData Map組裝數據
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{"name": deploymentName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}
	//序列化為字節，因為path方法只接收字節類型參數
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error("patchData序列化失敗", err)
		return errors.New("patchData序列化失敗" + err.Error())
	}
	//調用patch方法更新deployment副本數
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(),
		deploymentName,
		"application/strategic-merge-patch+json",
		patchByte,
		metav1.PatchOptions{})
	if err != nil {
		logger.Error("修改Deployment副本數失敗", err)
		return errors.New("修改Deployment副本數失敗" + err.Error())
	}
	return nil
}
//創建Deployment
//{
//  "name": "first-nginx",
//  "namespace": "default",
//  "replicas": 1,
//  "image": "nginx:latest",
//  "resource": "0.5/1",
//  "health_check": true,
//  "health_path": "/",
//  "label_str": "app=first-nginx",
//  "label": {
//    "app": "first-nginx"
//  },
//  "container_port": 80,
//  "cpu": "0.1",
//  "memory": "1Gi"
//}
func(p *deployment) CreateDeployment(data *DeployCreate) (err error) {
	//將data中的數據組裝成appsv1.Deployment對象
	deployment := &appsv1.Deployment{
		//元數據
		ObjectMeta: metav1.ObjectMeta{
			Name: data.Name,
			Namespace: data.Namespace,
			Labels: data.Label,
		},
		//副本數、選擇器，以及pod屬性
		Spec:       appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			//Pod數據
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: data.Name,
					Labels: data.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name: "http",
									Protocol: corev1.ProtocolTCP,
									ContainerPort: data.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		Status:     appsv1.DeploymentStatus{},
	}
	//判斷是否打開健康檢查功能，若打開，則則定ReadinessProbe和LivenessProbe
	if data.HealthCheck {
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			Handler:             corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:        data.HealthPath,
					Port:        intstr.IntOrString{
						Type: 0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			Handler:             corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:        data.HealthPath,
					Port:        intstr.IntOrString{
						Type: 0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}
	//定義容器的limit和request資源
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU: resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU: resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	//創建deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error("創建Deployment失敗", err)
		return errors.New("創建Deployment失敗" + err.Error())
	}

	return nil
}

//获取每个命名空间deployment的数量
func(p *deployment) GetDeploymentNumPerNs() (deploymentsNss []*DeploymentsNs, err error) {
	//获取namespace列表
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取Namespace列表失败", err)
		return nil, errors.New("获取Namespace列表失败" + err.Error())
	}
	//for循环
	for _, namespace := range namespaceList.Items {
		//获取deployment列表
		deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			//logger是给自己看的，return是给用户看的
			logger.Error("获取deployment列表失败", err)
			return nil, errors.New("获取deployment列表失败" + err.Error())
		}
		//组装数据
		deploymentsNs := &DeploymentsNs{
			Namespace: namespace.Name,
			DeploymentNum:    len(deploymentList.Items),
		}
		deploymentsNss = append(deploymentsNss, deploymentsNs)
	}
	return deploymentsNss, nil
}

//把deploymentCell转成appsv1 deployment
func(p *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}

//把appsv1 deployment转成DataCell
func(p *deployment) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i])
	}
	return cells
}