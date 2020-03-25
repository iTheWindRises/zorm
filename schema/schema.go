package schema

import (
	"go/ast"
	"reflect"
	"strconv"
	"strings"
	"zorm/dialect"
)

const (
	AUTO_INCREMENT = "auto_increment"
	PRIMARY_KEY = "primary key"
	NOT_NULL = "not null"
)

type Field struct {
	Name string	//字段名
	Type string	//类型
	Size int	//字段长度
	IsKey bool 	// 是否主键
	IsNotNil bool
	IsAutoIncrement bool	//是否自动递增
}

type Schema struct {
	Model reflect.Type	// 被映射对象
	Name string			// 表名
	Fields []*Field		// 字段
	FieldNames []string
	ColumnNames []string
	fieldMap map[string]*Field
	Key string 			// 主键名称
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}



func Parse(modelType reflect.Type, d dialect.Dialect) *Schema {
	//modelType := reflect.TypeOf(dest)
	schema := &Schema{
		Model:      modelType,
		Name:       strings.ToLower(modelType.Name()),
		fieldMap:   make(map[string]*Field),
	}

	for i:=0;i< modelType.NumField();i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			// 默认size
			if field.Type == "varchar" {
				field.Size = 255
			}
			if v, ok := p.Tag.Lookup("zorm"); ok {
				// 切割tag
				tags := strings.Split(v, ";")
				for _,tag := range tags {
					val := strings.ToLower(tag)
					if val == AUTO_INCREMENT {
						field.IsAutoIncrement = true
					}else if val == PRIMARY_KEY {
						field.IsKey = true
						schema.Key = field.Name
					}else if val == NOT_NULL {
						field.IsAutoIncrement  = true
					}else if strings.HasPrefix(val, "size:") {
						if size, err := strconv.Atoi(val[5:]); err == nil {
							field.Size = size
						}
					}else if strings.HasPrefix(val, "column") {
						field.Name = val[7:]
						if field.IsKey {
							schema.Key = field.Name
						}
					}else {

					}
				}
			}

			if schema.Key == "" {
				schema.Key = "id"
			}
			schema.Fields = append(schema.Fields,field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.ColumnNames = append(schema.ColumnNames,field.Name)

			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}