package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"k8s-platform/db"
	"k8s-platform/model"
)

var Workflow workflow

type workflow struct{}

//列表结构体
type WorkflowResp struct {
	Items []*model.Workflow `json:"items"`
	Total int  `json:"total"`
}

//列表
func(w *workflow) GetList(name string, page, limit int) (data *WorkflowResp, err error) {
	//定义分页数据的起始位置, 如果是第二页，每页10条数据, startSet=10
	startSet := (page-1) * limit
	//定义数据库查询返回内容
	var (
		workflowList []*model.Workflow
		total        int
	)
	tx := db.GORM.
		Model(&model.Workflow{}).  //指定表名,给count用的，否则无法找到对应的表
		Where("name like ?", "%" + name + "%").  //实现过滤
		Count(&total). //取总数
		Limit(limit).     //实现分页
		Offset(startSet).
		Order("id desc").   //实现排序
		Find(&workflowList)
	if tx.Error != nil {
		logger.Error("获取workflow列表失败", tx.Error)
		return nil, errors.New("获取workflow列表失败" + tx.Error.Error())
	}
	return &WorkflowResp{
		Items: workflowList,
		Total: total,
	}, nil
}

//获取单条
func(w *workflow) GetById(id int) (workflow *model.Workflow, err error) {
	//使用first或者find方法时，必须要初始化结构体（分配内存），不然会报错
	workflow = &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).Find(&workflow)
	if tx.Error != nil {
		logger.Error("获取workflow单条数据失败", tx.Error)
		return nil, errors.New("获取workflow单条数据失败" + tx.Error.Error())
	}
	return workflow, nil
}

//新增
func(w *workflow) Add(workflow *model.Workflow) (err error) {
	tx := db.GORM.Create(&workflow)
	if tx.Error != nil {
		logger.Error("创建workflow失败", tx.Error)
		return errors.New("创建workflow失败" + tx.Error.Error())
	}
	return nil
}
//删除
func(w *workflow) DelById(id int) (err error) {
	tx := db.GORM.Where("id = ?", id).Delete(&model.Workflow{})
	if tx.Error != nil {
		logger.Error("删除workflow失败", tx.Error)
		return errors.New("删除workflow失败" + tx.Error.Error())
	}
	return nil
}