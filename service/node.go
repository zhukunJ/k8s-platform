package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Node node

type node struct {}

//定义node列表的返回内容，Items是node列表，Total为node元素总数
//先过滤，再拿total，再做分页
type NodesResp struct{
	Items []corev1.Node `json:"items"`
	Total int         `json:"total"`
}

//获取node列表
func(p *node) GetNodes(filterName string, limit, page int) (nodesResp *NodesResp, err error) {
	//通过clintset获取node完整列表
	nodeList, err := K8s.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取node列表失败", err)
		return nil, errors.New("获取node列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(nodeList.Items),
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
	//再将DataCell切片数据转成原生node切片
	nodes := p.fromCells(data.GenericDataList)
	//返回
	return &NodesResp{
		Items: nodes,
		Total: total,
	}, nil
}
//获取node详情
func(p *node) GetNodeDetail(nodeName string) (node *corev1.Node, err error) {
	node, err = K8s.ClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Node详情失败", err)
		return nil, errors.New("获取Node详情失败" + err.Error())
	}
	return node, nil
}

//把nodeCell转成corev1 node
func(p *node) fromCells(cells []DataCell) []corev1.Node {
	nodes := make([]corev1.Node, len(cells))
	for i := range cells {
		nodes[i] = corev1.Node(cells[i].(nodeCell))
	}
	return nodes
}

//把corev1 node转成DataCell
func(p *node) toCells(std []corev1.Node) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = nodeCell(std[i])
	}
	return cells
}