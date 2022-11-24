package service

import (
	"errors"
	"k8s-platform/model"

	"github.com/wonderivan/logger"
)

var Login login

type login struct{}

//验证账号密码
func (l *login) Auth(username, password string) (udata *model.Iopsflow, err error) {
	data, err := Iopsflow.GetByName(username)
	if err != nil {
		logger.Error("获取用户信息失败,", err)
		return nil, err
	}

	// 如果密码不正确，返回错误
	if data.Password != password {
		logger.Error("登录失败, 用户名或密码错误")
		return nil, errors.New("登录失败, 用户名或密码错误")
	}
	// 如果密码正确，返回nil
	logger.Error("用户校验成功")
	return data, nil
}
