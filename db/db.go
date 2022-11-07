package db

import (
	"fmt"
	"k8s-platform/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //gorm对应的mysql驱动
	"github.com/wonderivan/logger"
)

var (
	isInit bool //是否已经初始化
	GORM   *gorm.DB
	err    error
)

//db的初始化函数，与数据库建立连接
func Init() {
	if isInit {
		return
	}
	//charset是数据库的字符集
	//parseTime是将数据库的时间类型自动转为go的时间类型
	//loc时区
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName)
	//建立连接
	GORM, err = gorm.Open(config.DbType, dsn)
	if err != nil {
		panic("数据库连接失败" + err.Error())
	}

	//设置debug开关
	GORM.LogMode(config.LogMode)

	//连接池配置
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	GORM.DB().SetConnMaxLifetime(config.MaxLifeTime)

	isInit = true
	logger.Info("数据库连接成功")
}
