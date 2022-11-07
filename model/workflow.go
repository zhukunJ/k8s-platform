package model

import "time"

//定义结构体，属性与mysql表字段对齐
type Workflow struct {
	//gorm:"primaryKey"用于声明主键
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Replicas   int32  `json:"replicas"`
	Deployment string `json:"deployment"`
	Service    string `json:"service"`
	Ingress    string `json:"ingress"`
	//gorm:"column:type"用于声明mysql中表的字段名
	Type string `json:"type" gorm:"column:type"`
}

//自定义设置表名
func (*Workflow) TableName() string {
	return "workflow"
}

//CREATE TABLE `workflow` (
// `id` int NOT NULL AUTO_INCREMENT,
// `name` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
// `namespace` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
// `replicas` int DEFAULT NULL,
// `deployment` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
// `service` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
// `ingress` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
// `type` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
// `created_at` datetime DEFAULT NULL,
// `updated_at` datetime DEFAULT NULL,
// `deleted_at` datetime DEFAULT NULL,
// PRIMARY KEY (`id`) USING BTREE,
// UNIQUE KEY `name` (`name`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
