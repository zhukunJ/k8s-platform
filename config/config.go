package config

import "time"

const (
	//gin监听的地址和端口
	ListenAddr = "0.0.0.0:9090"
	// windows
	// KubeConfig = "C:\\Users\\kunjiezx\\.kube\\config"
	// linux or mac
	KubeConfig = "/Users/admin/.kube/config"
	//查看日志的行数
	PodLogTailLine = 2000
	//管理员账号密码
	AdminUser = "admin"
	AdminPwd  = "123456"
	//数据库配置
	DbType     = "mysql"
	DbHost     = "127.0.0.1"
	DbPort     = 3306
	DbName     = "access"
	DbUser     = "root"
	DbPassword = "intel5g"
	//打印mysql debug日志开关
	LogMode = true
	//连接池配置
	MaxIdleConns = 10               //最大空闲连接
	MaxOpenConns = 100              //最大连接数
	MaxLifeTime  = 30 * time.Second //最大生存时间
)
