package session

import (
	"database/sql"
	"strings"

	"github.com/Muskchen/toolkits/gee/geeORM/clause"

	"github.com/Muskchen/toolkits/gee/geeORM/dialect"
	"github.com/Muskchen/toolkits/gee/geeORM/log"
	"github.com/Muskchen/toolkits/gee/geeORM/schema"
)

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	tx       *sql.Tx
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []interface{}
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var (
	_ CommonDB = (*sql.DB)(nil)
	_ CommonDB = (*sql.Tx)(nil)
)

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Raw(sql string, value ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, value...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return result, nil
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
