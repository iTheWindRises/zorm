package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"zorm/dialect"
	"zorm/log"
	"zorm/schema"
	"zorm/session"
)

// mysql session
type MSession struct {
	Db       *sql.DB
	sql      strings.Builder
	sqlVals  []interface{}
	refMap map[string]*schema.Schema
	dialect dialect.Dialect
	clause Clause
	tx *sql.Tx
}

func New(db *sql.DB, dialect dialect.Dialect) session.Session {
	return &MSession{
		Db:      db,
		dialect: dialect,
		refMap: make(map[string]*schema.Schema),
	}
}


func (s *MSession) DB() *sql.DB {
	return s.Db
}

func (s *MSession) Clear() {
	s.sql.Reset()
	s.sqlVals = nil
	s.clause = Clause{}
}

func (s *MSession) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVals)
	if result, err = s.DB().Exec(s.sql.String(),s.sqlVals...); err != nil {
		log.Error(err)
	}
	return
}

func (s *MSession) Raw(sql string, val ...interface{}) session.Session{
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVals = append(s.sqlVals, val...)
	return s
}

func (s *MSession) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVals)
	return s.DB().QueryRow(s.sql.String(),s.sqlVals...)
}

func (s *MSession) Query() *sql.Rows {
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVals)
	if rows, err := s.DB().Query(s.sql.String(),s.sqlVals...); err!= nil {
		log.Error(err.Error())
	}else {
		return rows
	}
	return nil
}


func (s *MSession) RefTable(key string) *schema.Schema {
	return s.refMap[key]
}

func (s *MSession) CreateTab(model interface{}) error {
	ts := s.Model(model)
	if s.HasTab(model) {
		msg := fmt.Sprintf("table %s exist.", s.RefTable(ts).Name)
		log.Info(msg)
		return errors.New(msg)
	}

	var columns[] string
	for _, field := range s.RefTable(ts).Fields {
		column := fmt.Sprintf("`%s` %s",field.Name,field.Type)
		if field.Size > 0 {
			column = fmt.Sprintf("%s(%d)",column,field.Size)
		}
		if field.IsNotNil {
			column = column +  " "+ schema.NOT_NULL
		}
		if field.IsAutoIncrement {
			column = column+ " " + schema.AUTO_INCREMENT
		}
		columns = append(columns,column)
	}
	desc := strings.Join(columns,",")
	_,err := s.Raw(fmt.Sprintf("CREATE TABLE `%s`(%s,PRIMARY KEY (`%s`))",
		s.RefTable(ts).Name, desc, s.RefTable(ts).Key)).Exec()

	return err
}
func (s *MSession) DropTab(model interface{}) error {
	_,err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s",s.RefTable("n").Name)).Exec()
	return err
}


func (s *MSession) Model(model interface{}) string {
	mtype := reflect.Indirect(reflect.ValueOf(model)).Type()
	if mtype.Kind() == reflect.Slice {
		mtype = mtype.Elem()
	}
	if _,ok := s.refMap[mtype.String()];!ok {
		s.refMap[mtype.String()] = schema.Parse(mtype,s.dialect)
	}
	return mtype.String()
}

func (s *MSession) HasTab(model interface{}) bool {
	ts := s.Model(model)
	sql := s.dialect.TableExistSQL(s.RefTable(ts).Name)
	row := s.Raw(sql).QueryRow()
	var tmp string
	if err := row.Scan(&tmp); err!= nil {
		log.Error(err.Error())
	}

	return tmp == s.RefTable(ts).Name
}

func (s *MSession) Insert(models ...interface{}) bool {
	ts := s.Model(models[0])
	s.clause.Set(INSERT,s.RefTable(ts).Name,s.RefTable(ts).ColumnNames)
	s.clause.Set(VALUES,getModelValues(models...)...)
	sql, vars := s.clause.Build(INSERT, VALUES)

	s.Raw(sql,vars...).Exec()
	return false
}


func (s *MSession) Find(m interface{}) error {
	// 获取m 的真实Value
	indirect := reflect.Indirect(reflect.ValueOf(m))
	// 获得m 的类型
	destType := indirect.Type()
	if destType.Kind() != reflect.Slice {
		return errors.New("params no slice.")
	}

	ts := s.Model(m)
	schema := s.RefTable(ts)
	s.clause.Set(SELECT,schema.Name,schema.ColumnNames)
	sql, vars := s.clause.Build(SELECT,WHERE,ORDERBY,LIMIT)
	s.Raw(sql,vars...)
	rows := s.Query()

	destType = destType.Elem()
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _,name := range schema.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		indirect.Set(reflect.Append(indirect,dest))
	}
	return nil
}

func (s *MSession) Where(sql string,vars ...interface{}) session.Session {
	s.clause.Set(WHERE,sql,vars)
	s.clause.Build(WHERE)
	return s
}

func (s *MSession) First(model interface{}) error {
	ts := s.Model(model)
	schema := s.RefTable(ts)
	s.clause.Set(SELECT,schema.Name,schema.ColumnNames)
	s.Limit(1)
	sql, vars := s.clause.Build(SELECT,WHERE,ORDERBY,LIMIT)
	s.Raw(sql,vars...)
	rows := s.Query()

	if rows.Next() {
		dest := reflect.ValueOf(model).Elem()
		var values []interface{}
		for _,name := range schema.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
	}
	return nil
}

func (s *MSession) OrderBy(s2 string) session.Session {
	s.clause.Set(ORDERBY,s2)
	s.clause.Build(ORDERBY)
	return s
}

func (s *MSession) Limit(limit int) session.Session {
	s.clause.Set(LIMIT,limit)
	s.clause.Build(LIMIT)
	return s
}

func (s *MSession) Save(model interface{}) error {
	ts := s.Model(model)
	schema := s.RefTable(ts)

	indirect := reflect.Indirect(reflect.ValueOf(model))
	var vars []interface{}
	keyVal := indirect.FieldByName(schema.FieldKey).Interface()
	for i:=0; i<indirect.NumField();i++ {
		vars = append(vars,indirect.Field(i).Interface())
	}
	s.clause.Set(UPDATE,schema.Name,schema.ColumnNames,vars)
	s.Where(fmt.Sprintf("%s=?",schema.Key),keyVal)
	sql, vars := s.clause.Build(UPDATE,WHERE)
	s.Raw(sql,vars...).Exec()
	return nil
}

func (s *MSession) Update(model interface{},pmap map[string]interface{}) error {
	ts := s.Model(model)
	schema := s.RefTable(ts)
	var cols []string
	var vars []interface{}
	for k,v := range pmap {
		cols = append(cols,k)
		vars = append(vars,v)
	}
	s.clause.Set(UPDATE,schema.Name,cols,vars)

	keyVal := reflect.Indirect(reflect.ValueOf(model)).FieldByName(schema.FieldKey).Interface()
	s.Where(fmt.Sprintf("%s=?",schema.Key),keyVal)
	sql, vars := s.clause.Build(UPDATE,WHERE)
	s.Raw(sql,vars...).Exec()
	return nil
}

func (s *MSession) Delete(model interface{}) error {
	ts := s.Model(model)
	schema := s.RefTable(ts)
	s.clause.Set(DELETE,schema.Name)
	s.Where(fmt.Sprintf("%s=?",schema.Key),reflect.Indirect(reflect.ValueOf(model)).FieldByName(schema.FieldKey).Interface())
	sql, vars := s.clause.Build(DELETE, WHERE)
	s.Raw(sql,vars...).Exec()
	return nil
}


func getModelValues(models ...interface{}) []interface{}  {
	var vals []interface{}
	for _, model := range models {
		var fieldValues []interface{}
		for i:=0;i< reflect.Indirect(reflect.ValueOf(model)).NumField();i++  {
			f := reflect.Indirect(reflect.ValueOf(model)).Field(i)
			fieldValues = append(fieldValues,f.Interface())
		}
		vals = append(vals,fieldValues)
	}
	return vals
}

