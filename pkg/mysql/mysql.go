package mysql

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	_defaultDriver          = "mysql"
	_defaultConnMaxLifetime = time.Minute * 3
	_defaultMaxOpenConns    = 20
	_defaultMaxIdleConns    = 20
)

type MySQL struct {
	Conn *sql.DB
}

func New(url string) (*MySQL, error) {
	md := &MySQL{}
	var err error
	md.Conn, err = sql.Open(_defaultDriver, url)

	if err != nil {
		return nil, fmt.Errorf("cannot connect to MySQL, error = %w", err)
	}

	md.Conn.SetConnMaxLifetime(_defaultConnMaxLifetime)
	md.Conn.SetMaxOpenConns(_defaultMaxOpenConns)
	md.Conn.SetMaxIdleConns(_defaultMaxIdleConns)

	if err = md.Conn.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping to MySQL, error = %w", err)
	}

	return md, nil
}
