package main

import (
	"fmt"
	"zorm"
	"zorm/model"
)
//CREATE TABLE `js_user` (
//	`uid` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
//	`username` varchar(31) NOT NULL DEFAULT '' COMMENT '用户名',
//	`nickname` varchar(31) NOT NULL DEFAULT '' COMMENT '昵称',
//	`gender` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '性别:0-男;1-女;2-保密',
//	`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
//	PRIMARY KEY (`uid`),
//	UNIQUE KEY `username` (`username`)
//)


func main() {
	engine, _ := zorm.NewEngine("mysql", "root:rootroot@tcp(116.62.57.22:3306)/test?charset=utf8mb4&loc=Local&parseTime=true")
	defer engine.Close()
	s := engine.NewSession()

	//s.Insert(u1,u2)
	var users []model.User
	s.Where("username = ?","张三").Find(&users)
	fmt.Println(users)
}
