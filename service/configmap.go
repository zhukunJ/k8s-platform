package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ConfigMap configMap

type configMap struct {}

//定义configMap列表的返回内容，Items是configMap列表，Total为configMap元素总数
//先过滤，再拿total，再做分页
type ConfigMapsResp struct{
	Items []corev1.ConfigMap `json:"items"`
	Total int         `json:"total"`
}

//获取configMap列表
func(p *configMap) GetConfigMaps(filterName, namespace string, limit, page int) (configMapsResp *ConfigMapsResp, err error) {
	//通过clintset获取configMap完整列表
	configMapList, err := K8s.ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取configMap列表失败", err)
		return nil, errors.New("获取configMap列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(configMapList.Items),
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
	//再将DataCell切片数据转成原生configMap切片
	configMaps := p.fromCells(data.GenericDataList)
	//返回
	return &ConfigMapsResp{
		Items: configMaps,
		Total: total,
	}, nil
}
//获取configMap详情
func(p *configMap) GetConfigMapDetail(configMapName, namespace string) (configMap *corev1.ConfigMap, err error) {
	configMap, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取ConfigMap详情失败", err)
		return nil, errors.New("获取ConfigMap详情失败" + err.Error())
	}
	return configMap, nil
}
//删除configMap
func(p *configMap) DeleteConfigMap(configMapName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), configMapName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除ConfigMap失败", err)
		return errors.New("删除ConfigMap失败" + err.Error())
	}
	return nil
}
//更新configMap
func(p *configMap) UpdateConfigMap(namespace, content string) (err error) {
	//将content反序列化成为configMap对象
	var configMap = &corev1.ConfigMap{}
	if err = json.Unmarshal([]byte(content), configMap); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新configMap
	_, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新ConfigMap失败", err)
		return errors.New("更新ConfigMap失败" + err.Error())
	}

	return nil
}

//把configMapCell转成corev1 configMap
func(p *configMap) fromCells(cells []DataCell) []corev1.ConfigMap {
	configMaps := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		configMaps[i] = corev1.ConfigMap(cells[i].(configMapCell))
	}
	return configMaps
}

//把corev1 configMap转成DataCell
func(p *configMap) toCells(std []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = configMapCell(std[i])
	}
	return cells
}