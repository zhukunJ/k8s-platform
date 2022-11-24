package model

import "time"

//定义结构体，属性与mysql表字段对齐
type Iopsflow struct {
	//gorm:"primaryKey"用于声明主键
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Username  string     `json:"username" gorm:"type:varchar(255);not null"`
	Password  string     `json:"password"`
	Editor    string     `json:"editor" gorm:"type:varchar(255);not null"`
	Read      bool       `json:"read"`
	Write     bool       `json:"write"`
	Delete    bool       `json:"delete"`
	Avatar    string     `json:"avatar"`
}

func (*Iopsflow) TableName() string {
	return "iops-access"
}

//数据库建表语句
// CREATE TABLE `iops-access` (
// 	`id` int NOT NULL AUTO_INCREMENT,
// 	`username` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
// 	`password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
//  `editor` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
// 	`read` boolean 		DEFAULT  0,
// 	`write` boolean 	DEFAULT 0,
// 	`delete` boolean 	DEFAULT 0,
// 	`avatar` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
// 	`created_at` datetime DEFAULT NULL,
// 	`updated_at` datetime DEFAULT NULL,
// 	`deleted_at` datetime DEFAULT NULL,
// 	PRIMARY KEY (`id`) USING BTREE,
// 	UNIQUE KEY `username` (`username`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci

// 增加数据
// insert into `iops-access` (`username`,`password`,`editor`,`read`,`write`,`delete`,`avatar`,`created_at`,`updated_at`,`deleted_at`) values ('intel','intel123','Admin',1,1,1,'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif',null,null,null)
