package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var StatefulSet statefulSet

type statefulSet struct {}

//定义statefulSet列表的返回内容，Items是statefulSet列表，Total为statefulSet元素总数
//先过滤，再拿total，再做分页
type StatefulSetsResp struct{
	Items []appsv1.StatefulSet `json:"items"`
	Total int         `json:"total"`
}

//获取statefulSet列表
func(p *statefulSet) GetStatefulSets(filterName, namespace string, limit, page int) (statefulSetsResp *StatefulSetsResp, err error) {
	//通过clintset获取statefulSet完整列表
	statefulSetList, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取statefulSet列表失败", err)
		return nil, errors.New("获取statefulSet列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(statefulSetList.Items),
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
	//再将DataCell切片数据转成原生statefulSet切片
	statefulSets := p.fromCells(data.GenericDataList)
	//返回
	return &StatefulSetsResp{
		Items: statefulSets,
		Total: total,
	}, nil
}
//获取statefulSet详情
func(p *statefulSet) GetStatefulSetDetail(statefulSetName, namespace string) (statefulSet *appsv1.StatefulSet, err error) {
	statefulSet, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取StatefulSet详情失败", err)
		return nil, errors.New("获取StatefulSet详情失败" + err.Error())
	}
	return statefulSet, nil
}
//删除statefulSet
func(p *statefulSet) DeleteStatefulSet(statefulSetName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除StatefulSet失败", err)
		return errors.New("删除StatefulSet失败" + err.Error())
	}
	return nil
}
//更新statefulSet
func(p *statefulSet) UpdateStatefulSet(namespace, content string) (err error) {
	//将content反序列化成为statefulSet对象
	var statefulSet = &appsv1.StatefulSet{}
	if err = json.Unmarshal([]byte(content), statefulSet); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新statefulSet
	_, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新StatefulSet失败", err)
		return errors.New("更新StatefulSet失败" + err.Error())
	}

	return nil
}

//把statefulSetCell转成appsv1 statefulSet
func(p *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulSets := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		statefulSets[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}
	return statefulSets
}

//把appsv1 statefulSet转成DataCell
func(p *statefulSet) toCells(std []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = statefulSetCell(std[i])
	}
	return cells
}