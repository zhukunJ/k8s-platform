package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	nwv1 "k8s.io/api/networking/v1"
	"sort"
	"strings"
	"time"
)

//定义DataSelector结构体，用于排序、过滤、分页
type DataSelector struct {
	GenericDataList []DataCell
	DataSelectQuery *DataSelect
}

//DataCell接口，用于各种资源的类型转换，排序、过滤、分页统一对DataCell进行处理
type DataCell interface {
	GetCreation() time.Time
	GetName()     string
}

//DataSelectQuery 定义过滤和分页的属性
type DataSelect struct {
	FilterQuery  *Filter
	PaginateQuery *Paginate
}

type Filter struct {
	Name string
}
type Paginate struct {
	Limit int
	Page  int
}

//排序
//实现自定义排序，需要重写Len、Swap、Less方法
//Len方法用于获取数组长度
func(d *DataSelector) Len() int {
	return len(d.GenericDataList)
}
//Swap方法用于在Less方法比较结果后，定义排序规则
func(d *DataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}
//Less方法用于定义数组中元素大小的比较方式
func(d *DataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}
//重写以上3个方法后，用sort.Sort方法进行排序
func(d *DataSelector) Sort() *DataSelector {
	sort.Sort(d)
	return d
}

//过滤
//比较元素中是否存在filterName相匹配的元素，若匹配，则返回
func(d *DataSelector) Filter() *DataSelector {
	//若Name传参为空，则返回所有
	if d.DataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	//若Name传参不为空，则返回切片中包含Name的所有元素
	filteredList := make([]DataCell,0)
	for _, value := range d.GenericDataList {
		matched := true
		objName := value.GetName()
		if !strings.Contains(objName, d.DataSelectQuery.FilterQuery.Name) {
			matched = false
			continue
		}
		if matched {
			filteredList = append(filteredList, value)
		}
	}
	d.GenericDataList = filteredList
	return d
}

//分页
//根据Limit和Page的传参，返回数据
func(d *DataSelector) Paginate() *DataSelector {
	limit := d.DataSelectQuery.PaginateQuery.Limit
	page := d.DataSelectQuery.PaginateQuery.Page
	//验证参数是否合法，若参数不合法，则返回所有
	if limit <= 0 || page <= 0 {
		return d
	}
	//举例：25个元素的数组，limit是10，page3 startIndex是20，endIndex25
	//第一页 0-10
	//第二页 10-20
	//第三页 20-25
	startIndex := limit * (page - 1)
	endIndex := limit * page
	//处理最后一页
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}

	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

//定义podCell类型，实现DataCell接口，用于类型转换
type podCell corev1.Pod

func(p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func(p podCell) GetName() string {
	return p.Name
}

//deployment
type deploymentCell appsv1.Deployment
func(d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}
func(d deploymentCell) GetName() string {
	return d.Name
}
//daemonset
type daemonSetCell appsv1.DaemonSet
func(d daemonSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}
func(d daemonSetCell) GetName() string {
	return d.Name
}
//statefulset
type statefulSetCell appsv1.StatefulSet
func(d statefulSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}
func(d statefulSetCell) GetName() string {
	return d.Name
}
//node
type nodeCell corev1.Node
func(p nodeCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func(p nodeCell) GetName() string {
	return p.Name
}
//namespace
type namespaceCell corev1.Namespace
func(p namespaceCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func(p namespaceCell) GetName() string {
	return p.Name
}
//pv
type pvCell corev1.PersistentVolume
func(p pvCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func(p pvCell) GetName() string {
	return p.Name
}
//service
type servicev1Cell corev1.Service
func(p servicev1Cell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func(p servicev1Cell) GetName() string {
	return p.Name
}
//ingress
type ingressCell nwv1.Ingress
func(p ingressCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func(p ingressCell) GetName() string {
	return p.Name
}
//configmap
type configMapCell corev1.ConfigMap
func(p configMapCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func(p configMapCell) GetName() string {
	return p.Name
}
//secret
type secretCell corev1.Secret
func(p secretCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func(p secretCell) GetName() string {
	return p.Name
}
//pvc
type pvcCell corev1.PersistentVolumeClaim
func(p pvcCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func(p pvcCell) GetName() string {
	return p.Name
}