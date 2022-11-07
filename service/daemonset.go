package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DaemonSet daemonSet

type daemonSet struct {}

//定义daemonSet列表的返回内容，Items是daemonSet列表，Total为daemonSet元素总数
//先过滤，再拿total，再做分页
type DaemonSetsResp struct{
	Items []appsv1.DaemonSet `json:"items"`
	Total int         `json:"total"`
}

//获取daemonSet列表
func(p *daemonSet) GetDaemonSets(filterName, namespace string, limit, page int) (daemonSetsResp *DaemonSetsResp, err error) {
	//通过clintset获取daemonSet完整列表
	daemonSetList, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取daemonSet列表失败", err)
		return nil, errors.New("获取daemonSet列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(daemonSetList.Items),
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
	//再将DataCell切片数据转成原生daemonSet切片
	daemonSets := p.fromCells(data.GenericDataList)
	//返回
	return &DaemonSetsResp{
		Items: daemonSets,
		Total: total,
	}, nil
}
//获取daemonSet详情
func(p *daemonSet) GetDaemonSetDetail(daemonSetName, namespace string) (daemonSet *appsv1.DaemonSet, err error) {
	daemonSet, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取DaemonSet详情失败", err)
		return nil, errors.New("获取DaemonSet详情失败" + err.Error())
	}
	return daemonSet, nil
}
//删除daemonSet
func(p *daemonSet) DeleteDaemonSet(daemonSetName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除DaemonSet失败", err)
		return errors.New("删除DaemonSet失败" + err.Error())
	}
	return nil
}
//更新daemonSet
func(p *daemonSet) UpdateDaemonSet(namespace, content string) (err error) {
	//将content反序列化成为daemonSet对象
	var daemonSet = &appsv1.DaemonSet{}
	if err = json.Unmarshal([]byte(content), daemonSet); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新daemonSet
	_, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新DaemonSet失败", err)
		return errors.New("更新DaemonSet失败" + err.Error())
	}

	return nil
}

//把daemonSetCell转成appsv1 daemonSet
func(p *daemonSet) fromCells(cells []DataCell) []appsv1.DaemonSet {
	daemonSets := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		daemonSets[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}
	return daemonSets
}

//把appsv1 daemonSet转成DataCell
func(p *daemonSet) toCells(std []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = daemonSetCell(std[i])
	}
	return cells
}