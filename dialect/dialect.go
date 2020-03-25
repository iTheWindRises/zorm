package dialect

import (
	"reflect"
)

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string	//用于将GO类型转化为该数据库类型
	TableExistSQL(tableName string)(string)	// 返回某表是否存在的sql语句
	TypeName() string
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name]=dialect
}

func GetDialect(name string)(dialect Dialect,ok bool)  {
	dialect, ok = dialectsMap[name]
	return
}