package session

import (
	"database/sql"
	"zorm/schema"
)

type Session interface {
	DB() *sql.DB
	Clear()
	Exec() (result sql.Result, err error)
	Raw(sql string, val ...interface{}) Session
	QueryRow() *sql.Row
	RefTable(ts string) *schema.Schema		//获得schema

	// ddl
	CreateTab(model interface{}) error
	DropTab(model interface{}) error
	Model(model interface{}) string
	HasTab(model interface{}) bool

	// dml
	Insert(model ...interface{}) bool
	Find(model interface{}) error
	Where(string,...interface{}) Session
}





