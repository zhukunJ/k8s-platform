package dao

import (
	"errors"
	"k8s-platform/db"
	"k8s-platform/model"

	"github.com/wonderivan/logger"
)

var Iopsflow iopsflow

type iopsflow struct{}

type IopsflowResp struct {
	Items []*model.Iopsflow
	Total int
}

//获取iopsflow列表
func (w *iopsflow) GetIopsflows(filterName string, limit, page int) (data *IopsflowResp, err error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	//定义数据库查询返回的内容
	var (
		iopsflowList []*model.Iopsflow
		total        int
	)
	//数据库查询，Limit方法用于限制条数，Offset方法用于设置起始位置
	tx := db.GORM.
		Model(&model.Iopsflow{}).
		Where("username like ?", "%"+filterName+"%").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&iopsflowList)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Iopsflow列表失败, " + tx.Error.Error())
		return nil, errors.New("获取Iopsflow列表失败, " + tx.Error.Error())
	}
	return &IopsflowResp{
		Items: iopsflowList,
		Total: total,
	}, nil
}

//获取详情
func (w *iopsflow) GetByName(username string) (iopsflow *model.Iopsflow, err error) {
	iopsflow = &model.Iopsflow{}
	tx := db.GORM.Where("username = ?", username).First(&iopsflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Iopsflow详情失败, " + tx.Error.Error())
		return nil, errors.New("获取Iopsflow详情失败, " + tx.Error.Error())
	}
	return iopsflow, nil
}

// //创建
func (w *iopsflow) Add(iopsflow *model.Iopsflow) (err error) {
	tx := db.GORM.Create(&iopsflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("创建Iopsflow失败, " + tx.Error.Error())
		return errors.New("创建Iopsflow失败, " + tx.Error.Error())
	}
	return nil
}

// //删除
func (w *iopsflow) DelByName(username string) (err error) {
	tx := db.GORM.Where("username = ?", username).Delete(&model.Iopsflow{})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Iopsflow详情失败, " + tx.Error.Error())
		return errors.New("获取Iopsflow详情失败, " + tx.Error.Error())
	}
	return nil
}
