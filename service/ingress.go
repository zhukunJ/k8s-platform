package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Ingress ingress

type ingress struct {}

//定义ingress列表的返回内容，Items是ingress列表，Total为ingress元素总数
//先过滤，再拿total，再做分页
type IngresssResp struct{
	Items []nwv1.Ingress `json:"items"`
	Total int         `json:"total"`
}

//定义ServiceCreate结构体，用于创建service需要的参数属性的定义
type IngressCreate struct {
	Name         string  `json:"name"`
	Namespace    string  `json:"namespace"`
	Label        map[string]string  `json:"label"`
	Hosts        map[string][]*HttpPath `json:"hosts"`
}
//定义ingress的path结构体
type HttpPath struct {
	Path         string         `json:"path"`
	PathType     nwv1.PathType  `json:"path_type"`
	ServiceName  string         `json:"service_name"`
	ServicePort  int32          `json:"service_port"`
}

//获取ingress列表
func(p *ingress) GetIngresss(filterName, namespace string, limit, page int) (ingresssResp *IngresssResp, err error) {
	//通过clintset获取ingress完整列表
	ingressList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取ingress列表失败", err)
		return nil, errors.New("获取ingress列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(ingressList.Items),
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
	//再将DataCell切片数据转成原生ingress切片
	ingresss := p.fromCells(data.GenericDataList)
	//返回
	return &IngresssResp{
		Items: ingresss,
		Total: total,
	}, nil
}
//获取ingress详情
func(p *ingress) GetIngressDetail(ingressName, namespace string) (ingress *nwv1.Ingress, err error) {
	ingress, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Ingress详情失败", err)
		return nil, errors.New("获取Ingress详情失败" + err.Error())
	}
	return ingress, nil
}
//创建ingress
//{
//  "name": "tst-ing1",
//  "namespace": "default",
//  "label_str": "app=first-nginx",
//  "host": "aaa.abc.com",
//  "path": "/",
//  "path_type": "Prefix",
//  "hosts": {
//    "aaa.abc.com": [
//      {
//        "path": "/",
//        "path_type": "Prefix",
//        "service_name": "tst-ing1",
//        "service_port": 80
//      }
//    ]
//  },
//  "label": {
//    "app": "first-nginx"
//  }
//}
func(i *ingress) CreateIngress(data *IngressCreate) (err error) {
	//声明nwv1.IngressRule和nwv1.HTTPIngressPath变量，后面组装数据于鏊用到
	var ingressRules []nwv1.IngressRule
	var httpIngressPATHs []nwv1.HTTPIngressPath
	//将data中的数据组装成nwv1.Ingress对象
	ingress := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: data.Name,
			Namespace: data.Namespace,
			Labels: data.Label,
		},
		Status: nwv1.IngressStatus{},
	}
	//第一层for循环是将host组装成nwv1.IngressRule类型的对象
	// 一个host对应一个ingressrule，每个ingressrule中包含一个host和多个path
	for key, value := range data.Hosts {
		ir := nwv1.IngressRule{
			Host: key,
			//这里现将nwv1.HTTPIngressRuleValue类型中的Paths置为空，后面组装好数据再赋值
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{Paths:nil},
			},
		}
		//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, httpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path: httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: getServiceName(httpPath.ServiceName),
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}
			//将每个hip对象组装成数组
			httpIngressPATHs = append(httpIngressPATHs, hip)
		}
		//给Paths赋值，前面置为空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPATHs
		//将每个ir对象组装成数组，这个ir对象就是IngressRule，每个元素是一个host和多个path
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules对象加入到ingress的规则中
	ingress.Spec.Rules = ingressRules
	//创建ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Ingress失败, " + err.Error()))
		return errors.New("创建Ingress失败, " + err.Error())
	}

	return nil
}
//删除ingress
func(p *ingress) DeleteIngress(ingressName, namespace string) (err error) {
	err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Ingress失败", err)
		return errors.New("删除Ingress失败" + err.Error())
	}
	return nil
}
//更新ingress
func(p *ingress) UpdateIngress(namespace, content string) (err error) {
	//将content反序列化成为ingress对象
	var ingress = &nwv1.Ingress{}
	if err = json.Unmarshal([]byte(content), ingress); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Ingress失败", err)
		return errors.New("更新Ingress失败" + err.Error())
	}

	return nil
}

//把ingressCell转成nwv1 ingress
func(p *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingresss := make([]nwv1.Ingress, len(cells))
	for i := range cells {
		ingresss[i] = nwv1.Ingress(cells[i].(ingressCell))
	}
	return ingresss
}

//把nwv1 ingress转成DataCell
func(p *ingress) toCells(std []nwv1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = ingressCell(std[i])
	}
	return cells
}