package main

import (
	"flag"
	"fmt"
	"time"
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

var (
	username string
	password string
	addr string
	database string
)

func init()  {
	flag.StringVar(&username,"u","root","数据库用户名")
	flag.StringVar(&password,"p","root","数据库用户密码")
	flag.StringVar(&addr,"addr","localhost:3306","数据库地址")
	flag.StringVar(&database,"use","test","数据库名称")

}

func main() {
	flag.Parse()
	engine, _ := zorm.NewEngine("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Local&parseTime=true",username,password,addr,database))
	defer engine.Close()
	s := engine.NewSession()

	//s.Insert(u1,u2)
	//var users []model.User
	//s.Limit(2).Where("username like  ?","%1%").OrderBy("create_at desc").Find(&users)
	//fmt.Println(users)

	//var user model.User
	//s.Where("username like  ?","%1%").First(&user)
	//fmt.Println(user)
	
	// update
	u1 := &model.User{
		Uid:       4,
		Username:  "爽歪歪",
		Nickname:  "故事2",
		Gender:    1,
		CreatedAt: time.Now(),
	}
	//s.Save(u1)
	//s.Update(u1, map[string]interface{}{"nickname":"故事3"})

	s.Delete(u1)
}
