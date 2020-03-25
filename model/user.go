package model

import "time"

type User struct {
	Uid int `zorm:"auto_increment;primary key;not null;size:11;column:uid"`
	Username string `zorm:"size:36;column:username"`
	Nickname string `zorm:"size:31;column:nickname"`
	Gender int `zorm:"size:3;column:gender"`
	CreatedAt time.Time `zorm:"column:create_at"`
}
