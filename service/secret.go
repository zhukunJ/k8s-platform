package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Secret secret

type secret struct {}

//定义secret列表的返回内容，Items是secret列表，Total为secret元素总数
//先过滤，再拿total，再做分页
type SecretsResp struct{
	Items []corev1.Secret `json:"items"`
	Total int         `json:"total"`
}

//获取secret列表
func(p *secret) GetSecrets(filterName, namespace string, limit, page int) (secretsResp *SecretsResp, err error) {
	//通过clintset获取secret完整列表
	secretList, err := K8s.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger是给自己看的，return是给用户看的
		logger.Error("获取secret列表失败", err)
		return nil, errors.New("获取secret列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(secretList.Items),
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
	//再将DataCell切片数据转成原生secret切片
	secrets := p.fromCells(data.GenericDataList)
	//返回
	return &SecretsResp{
		Items: secrets,
		Total: total,
	}, nil
}
//获取secret详情
func(p *secret) GetSecretDetail(secretName, namespace string) (secret *corev1.Secret, err error) {
	secret, err = K8s.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Secret详情失败", err)
		return nil, errors.New("获取Secret详情失败" + err.Error())
	}
	return secret, nil
}
//删除secret
func(p *secret) DeleteSecret(secretName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Secret失败", err)
		return errors.New("删除Secret失败" + err.Error())
	}
	return nil
}
//更新secret
func(p *secret) UpdateSecret(namespace, content string) (err error) {
	//将content反序列化成为secret对象
	var secret = &corev1.Secret{}
	if err = json.Unmarshal([]byte(content), secret); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新secret
	_, err = K8s.ClientSet.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Secret失败", err)
		return errors.New("更新Secret失败" + err.Error())
	}

	return nil
}

//把secretCell转成corev1 secret
func(p *secret) fromCells(cells []DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		secrets[i] = corev1.Secret(cells[i].(secretCell))
	}
	return secrets
}

//把corev1 secret转成DataCell
func(p *secret) toCells(std []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = secretCell(std[i])
	}
	return cells
}