package zorm

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"zorm/dialect"
	"zorm/log"
	"zorm/mysql"
	"zorm/session"
)

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine,err error) {
	dialect,err := selectParser(driver)
	if err != nil {
		log.Error(err)
		return
	}
	db, err := sql.Open(driver,source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	e = &Engine{
		db: db,
		dialect:dialect,
	}
	log.Info("Connect database success")
	return
}

func (e *Engine)Close()  {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e *Engine)NewSession() session.Session {
	return selectSession(e.db,e.dialect)
}

func selectParser(driver string) (d dialect.Dialect,err error) {
	if d, ok := dialect.GetDialect(driver); ok {
		return d,nil
	}else {
		err = errors.New("dialect parse fail.")
		return nil,err
	}
}

// 选择数据解析器
func selectSession(db *sql.DB,dialect dialect.Dialect) session.Session {
	switch dialect.TypeName() {
	case "mysql":
		return mysql.New(db,dialect)
	}
	return nil
}