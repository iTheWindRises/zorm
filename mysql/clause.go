package mysql

import (
	"fmt"
	"strings"
)

type Type int
const (
	INSERT Type = iota
	VALUES
	SELECT
	WHERE
	ORDERBY
	LIMIT
	UPDATE
	DELETE
	COUNT
)

type Clause struct {
	sql map[Type]string
	sqlVals map[Type][]interface{}
}

func (c *Clause)Set(name Type, vars ...interface{})  {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVals = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVals[name] = vars
}

func (c *Clause)Build(orders ...Type)(string, []interface{})  {
	var sqls []string
	var vars []interface{}
	for _,order := range orders {
		if sql, ok := c.sql[order];ok {
			sqls = append(sqls,sql)
			vars = append(vars,c.sqlVals[order]...)
		}
	}
	return strings.Join(sqls," "),vars
}


type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator



func init()  {
	generators = make(map[Type]generator)
	generators[INSERT] =_insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tablename ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string),",")
	return fmt.Sprintf("INSERT INTO %s (%v)",tableName,fields),[]interface{}{}
}


func _update(values ...interface{}) (string, []interface{}) {
	// UPDATE $tableName SET $field=?,... WHERE $condition=?            INTO $tablename ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string),"=?,")+ "=?"
	values = values[2:]
	return fmt.Sprintf("UPDATE `%s` SET %v",tableName,fields),values[0].([]interface{})
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _values(values ...interface{}) (string, []interface{}) {
	// VALUES($v1),($v2)
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, val := range values {
		v := val.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)",bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars,v...)
	}
	return sql.String(),vars
}

func _select(values ...interface{}) (string, []interface{})  {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string),",")
	return fmt.Sprintf("SELECT %v FROM %s",fields,tableName),[]interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars[0].([]interface{})
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

func _delete(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	tab := values[0]
	return fmt.Sprintf("DELETE FROM `%s`", tab), []interface{}{}
}

