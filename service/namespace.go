package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Namespace namespace

type namespace struct {}

//定义namespace列表的返回内容，Items是namespace列表，Total为namespace元素总数
//先过滤，再拿total，再做分页
type NamespacesResp struct{
	Items []corev1.Namespace `json:"items"`
	Total int         `json:"total"`
}

//获取namespace列表
func(p *namespace) GetNamespaces(filterName string, limit, page int) (namespacesResp *NamespacesResp, err error) {
	//通过clintset获取namespace完整列表
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取namespace列表失败", err)
		return nil, errors.New("获取namespace列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(namespaceList.Items),
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
	//再将DataCell切片数据转成原生namespace切片
	namespaces := p.fromCells(data.GenericDataList)
	//返回
	return &NamespacesResp{
		Items: namespaces,
		Total: total,
	}, nil
}
//获取namespace详情
func(p *namespace) GetNamespaceDetail(namespaceName string) (namespace *corev1.Namespace, err error) {
	namespace, err = K8s.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Namespace详情失败", err)
		return nil, errors.New("获取Namespace详情失败" + err.Error())
	}
	return namespace, nil
}
//删除namespace
func(p *namespace) DeleteNamespace(namespaceName string) (err error) {
	err = K8s.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Namespace失败", err)
		return errors.New("删除Namespace失败" + err.Error())
	}
	return nil
}

//把namespaceCell转成corev1 namespace
func(p *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	namespaces := make([]corev1.Namespace, len(cells))
	for i := range cells {
		namespaces[i] = corev1.Namespace(cells[i].(namespaceCell))
	}
	return namespaces
}

//把corev1 namespace转成DataCell
func(p *namespace) toCells(std []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = namespaceCell(std[i])
	}
	return cells
}