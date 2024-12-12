package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/idoyudha/eshop-cart/config"
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

func NewMySQL(cfg config.MySQL) (*MySQL, error) {
	mysql := &MySQL{}
	var err error

	mysql.Conn, err = sql.Open(_defaultDriver, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	mysql.Conn.SetConnMaxLifetime(time.Minute * time.Duration(cfg.ConnectionMaxLifetime))
	mysql.Conn.SetMaxOpenConns(cfg.MaxOpenConnection)
	mysql.Conn.SetMaxIdleConns(cfg.MaxIdleConnection)

	if err = mysql.Conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return mysql, nil
}
