package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var Servicev1 servicev1

type servicev1 struct {}

//定义servicev1列表的返回内容，Items是servicev1列表，Total为servicev1元素总数
//先过滤，再拿total，再做分页
type Servicev1sResp struct{
	Items []corev1.Service `json:"items"`
	Total int         `json:"total"`
}

type ServiceCreate struct {
	Name           string  `json:"name"`
	Namespace      string  `json:"namespace"`
	Type           string  `json:"type"`
	ContainerPort  int32   `json:"container_port"`
	Port           int32   `json:"port"`
	NodePort       int32   `json:"node_port"`
	Label          map[string]string  `json:"label"`
}

//获取servicev1列表
func(p *servicev1) GetServicev1s(filterName, namespace string, limit, page int) (servicev1sResp *Servicev1sResp, err error) {
	//通过clintset获取servicev1完整列表
	servicev1List, err := K8s.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取servicev1列表失败", err)
		return nil, errors.New("获取servicev1列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(servicev1List.Items),
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
	//再将DataCell切片数据转成原生servicev1切片
	servicev1s := p.fromCells(data.GenericDataList)
	//返回
	return &Servicev1sResp{
		Items: servicev1s,
		Total: total,
	}, nil
}
//获取servicev1详情
func(p *servicev1) GetServicev1Detail(servicev1Name, namespace string) (servicev1 *corev1.Service, err error) {
	servicev1, err = K8s.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), servicev1Name, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Servicev1详情失败", err)
		return nil, errors.New("获取Servicev1详情失败" + err.Error())
	}
	return servicev1, nil
}
//创建service,,接收ServiceCreate对象
//{
//  "name": "t-svc1",
//  "namespace": "default",
//  "type": "ClusterIP",
//  "container_port": 80,
//  "port": 80,
//  "label": {
//    "app": "first-nginx"
//  }
//}
func(s *servicev1) CreateServicev1(data *ServiceCreate) (err error) {
	//将data中的数据组装成corev1.Service对象
	service := &corev1.Service{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name: data.Name,
			Namespace: data.Namespace,
			Labels: data.Label,
		},
		//Spec中定义类型，端口，选择器
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: data.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type: 0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Label,
		},
	}
	//默认ClusterIP,这里是判断NodePort,添加配置
	if data.NodePort != 0 && data.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = data.NodePort
	}
	//创建Service
	_, err = K8s.ClientSet.CoreV1().Services(data.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Service失败, " + err.Error()))
		return errors.New("创建Service失败, " + err.Error())
	}

	return nil
}
//删除servicev1
func(p *servicev1) DeleteServicev1(servicev1Name, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), servicev1Name, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Servicev1失败", err)
		return errors.New("删除Servicev1失败" + err.Error())
	}
	return nil
}
//更新servicev1
func(p *servicev1) UpdateServicev1(namespace, content string) (err error) {
	//将content反序列化成为servicev1对象
	var servicev1 = &corev1.Service{}
	if err = json.Unmarshal([]byte(content), servicev1); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新servicev1
	_, err = K8s.ClientSet.CoreV1().Services(namespace).Update(context.TODO(), servicev1, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Servicev1失败", err)
		return errors.New("更新Servicev1失败" + err.Error())
	}

	return nil
}

//把servicev1Cell转成corev1 servicev1
func(p *servicev1) fromCells(cells []DataCell) []corev1.Service {
	servicev1s := make([]corev1.Service, len(cells))
	for i := range cells {
		servicev1s[i] = corev1.Service(cells[i].(servicev1Cell))
	}
	return servicev1s
}

//把corev1 servicev1转成DataCell
func(p *servicev1) toCells(std []corev1.Service) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = servicev1Cell(std[i])
	}
	return cells
}