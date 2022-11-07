package service

import (
	"k8s-platform/dao"
	"k8s-platform/model"
)

var Iopsflow iopsflow

type iopsflow struct{}

//定义iopsflowCreate类型
type IopsflowCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Editor   string `json:"editor"`
	Read     bool   `json:"read"`
	Write    bool   `json:"write"`
	Delete   bool   `json:"delete"`
	Avatar   string `json:"avatar"`
}

//获取列表分页查询
func (w *iopsflow) GetList(name string, limit, page int) (data *dao.IopsflowResp, err error) {
	data, err = dao.Iopsflow.GetIopsflows(name, limit, page)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//查询iopsflow单条数据
func (w *iopsflow) GetByName(username string) (data *model.Iopsflow, err error) {
	data, err = dao.Iopsflow.GetByName(username)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// //创建iopsflow
func (w *iopsflow) CreateWorkFlow(data *IopsflowCreate) (err error) {
	//iopsflow数据落库
	iopsflow := &model.Iopsflow{
		Username: data.Username,
		Password: data.Password,
		Editor:   data.Editor,
		Read:     data.Read,
		Write:    data.Write,
		Delete:   data.Delete,
		Avatar:   data.Avatar,
	}

	err = dao.Iopsflow.Add(iopsflow)
	if err != nil {
		return err
	}
	return err
}

// //删除iopsflow
func (w *iopsflow) DelByName(username string) (err error) {

	//删除数据库数据
	err = dao.Iopsflow.DelByName(username)
	if err != nil {
		return err
	}

	return
}

// //删除k8s资源 deployment service ingress
// func delIopsflowRes(iopsflow *model.Iopsflow) (err error) {
// 	err = Deployment.DeleteDeployment(iopsflow.Name, iopsflow.Namespace)
// 	if err != nil {
// 		return err
// 	}
// 	err = Servicev1.DeleteService(getServiceName(iopsflow.Name), iopsflow.Namespace)
// 	if err != nil {
// 		return err
// 	}

// 	if iopsflow.Type == "Ingress" {
// 		err = Ingress.DeleteIngress(getIngressName(iopsflow.Name), iopsflow.Namespace)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// //创建k8s资源 deployment service ingress
// func createIopsflowRes(data *IopsflowCreate) (err error) {

// 	//创建deployment
// 	dc := &DeployCreate{
// 		Name:          data.Name,
// 		Namespace:     data.Namespace,
// 		Replicas:      data.Replicas,
// 		Image:         data.Image,
// 		Label:         data.Label,
// 		Cpu:           data.Cpu,
// 		Memory:        data.Memory,
// 		ContainerPort: data.ContainerPort,
// 		HealthCheck:   data.HealthCheck,
// 		HealthPath:    data.HealthPath,
// 	}
// 	err = Deployment.CreateDeployment(dc)
// 	if err != nil {
// 		return err
// 	}
// 	var serviceType string
// 	if data.Type != "Ingress" {
// 		serviceType = data.Type
// 	} else {
// 		serviceType = "ClusterIP"
// 	}
// 	//创建service
// 	sc := &ServiceCreate{
// 		Name:          getServiceName(data.Name),
// 		Namespace:     data.Namespace,
// 		Type:          serviceType,
// 		ContainerPort: data.ContainerPort,
// 		Port:          data.Port,
// 		NodePort:      data.NodePort,
// 		Label:         data.Label,
// 	}
// 	if err := Servicev1.CreateService(sc); err != nil {
// 		return err
// 	}
// 	//创建ingress
// 	var ic *IngressCreate
// 	if data.Type == "Ingress" {
// 		ic = &IngressCreate{
// 			Name:      getIngressName(data.Name),
// 			Namespace: data.Namespace,
// 			Label:     data.Label,
// 			Hosts:     data.Hosts,
// 		}
// 	}
// 	err = Ingress.CreateIngress(ic)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// //iopsflow名字转换成service名字，添加-svc后缀
// func getServiceName(iopsflowName string) (serviceName string) {
// 	return iopsflowName + "-svc"
// }

// //iopsflow名字转换成ingress名字，添加-ing后缀
// func getIngressName(iopsflowName string) (ingressName string) {
// 	return iopsflowName + "-ing"
// }
