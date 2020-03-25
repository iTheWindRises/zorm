package mysql

import (
	"fmt"
	"reflect"
	"time"
	"zorm/dialect"
)
func init()  {
	dialect.RegisterDialect("mysql",&mparser{})
}
type mparser struct {
}

func (m *mparser) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "tinyint"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "int"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "varchar"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "timestamp"
		}
	case reflect.Ptr:
		if typ.Type().String() == "*time.Time" {
			return "timestamp"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (m *mparser) TableExistSQL(tableName string) string {
	return "SHOW TABLES LIKE '"+tableName+"'"
}

func (m *mparser) TypeName() string {
	return "mysql"
}
