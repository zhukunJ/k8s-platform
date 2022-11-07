package service

import (
	"k8s-platform/dao"
	"k8s-platform/model"
)

var Workflow workflow

type workflow struct {}

//定义WorkflowCreate结构体，用于创建workflow需要的参数属性的定义
type WorkflowCreate struct {
	Name           string  `json:"name"`
	Namespace      string  `json:"namespace"`
	Replicas       int32   `json:"replicas"`
	Image          string  `json:"image"`
	Label          map[string]string  `json:"label"`
	Cpu            string  `json:"cpu"`
	Memory         string  `json:"memory"`
	ContainerPort  int32   `json:"container_port"`
	HealthCheck    bool    `json:"health_check"`
	HealthPath     string  `json:"health_path"`
	Type           string  `json:"type"`
	Port           int32   `json:"port"`
	NodePort       int32   `json:"node_port"`
	Hosts          map[string][]*HttpPath `json:"hosts"`
}

//获取列表分页查询
func(w *workflow) GetList(name string, page, limit int) (data *dao.WorkflowResp, err error) {
	data, err = dao.Workflow.GetList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//查询workflow单条数据
func(w *workflow) GetById(id int) (data *model.Workflow, err error) {
	data, err = dao.Workflow.GetById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//新增
func(w *workflow) CreateWorkflow(data *WorkflowCreate) (err error) {
	//处理数据
	var ingressName string
	if data.Type == "Ingress" {
		ingressName = getIngressName(data.Name)
	} else {
		ingressName = ""
	}
	workflow := &model.Workflow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicas,
		Deployment: data.Name,
		Service:    getServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
	}
	//创建数据库记录
	err = dao.Workflow.Add(workflow)
	if err != nil {
		return err
	}
	//创建k8s资源
	err = createWorkflowRes(data)
	if err != nil {
		return err
	}
	return nil
}
//删除workflow
func(w *workflow) DelById(id int) (err error) {
	//获取workflow数据
	workflow, err := dao.Workflow.GetById(id)
	if err != nil {
		return err
	}
	//删除k8s资源
	err = delWorkflowRes(workflow)
	if err != nil {
		return err
	}
	//删除数据库数据
	err = dao.Workflow.DelById(id)
	if err != nil {
		return err
	}

	return nil
}
//新增，创建k8s资源
func createWorkflowRes(data *WorkflowCreate) (err error) {
	//创建deployment
	dc := &DeployCreate{
		Name:          data.Name,
		Namespace:     data.Namespace,
		Replicas:      data.Replicas,
		Image:         data.Image,
		Label:         data.Label,
		Cpu:           data.Cpu,
		Memory:        data.Memory,
		ContainerPort: data.ContainerPort,
		HealthCheck:   data.HealthCheck,
		HealthPath:    data.HealthPath,
	}
	err = Deployment.CreateDeployment(dc)
	if err != nil {
		return err
	}
	//判断类型
	var serviceType string
	if data.Type != "Ingress" {
		serviceType = data.Type
	} else {
		serviceType = "ClusterIP"
	}
	//创建service
	sc := &ServiceCreate{
		Name:          getServiceName(data.Name),
		Namespace:     data.Namespace,
		Type:          serviceType,
		ContainerPort: data.ContainerPort,
		Port:          data.Port,
		NodePort:      data.NodePort,
		Label:         data.Label,
	}
	err = Servicev1.CreateServicev1(sc)
	if err != nil {
		return err
	}
	//创建ingress
	if data.Type == "Ingress" {
		ic := &IngressCreate{
			Name:      getIngressName(data.Name),
			Namespace: data.Namespace,
			Label:     data.Label,
			Hosts:     data.Hosts,
		}
		err = Ingress.CreateIngress(ic)
		if err != nil {
			return err
		}
	}
	return nil
}
//封装删除workflow对应的k8s资源
func delWorkflowRes(workflow *model.Workflow) (err error) {
	//删除deployment
	err = Deployment.DeleteDeployment(workflow.Name, workflow.Namespace)
	if err != nil {
		return err
	}
	//删除service
	err = Servicev1.DeleteServicev1(getServiceName(workflow.Name), workflow.Namespace)
	if err != nil {
		return err
	}
	//删除ingress，这里多了一层判断，因为只有type为ingress的workflow才有ingress资源
	if workflow.Type == "Ingress" {
		err = Ingress.DeleteIngress(getIngressName(workflow.Name), workflow.Namespace)
		if err != nil {
			return err
		}
	}

	return nil
}
//workflow名字转换成service名字，添加-svc后缀
func getServiceName(workflowName string) (serviceName string) {
	return workflowName + "-svc"
}
//workflow名字转换成ingress名字，添加-ing后缀
func getIngressName(workflowName string) (ingressName string) {
	return workflowName + "-ing"
}