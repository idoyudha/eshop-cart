package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/idoyudha/eshop-cart/config"
)

const (
	_defaultDriver          = "mysql"
	_defaultConnMaxLifetime = 2 * time.Second
	_defaultMaxOpenConns    = 4 // (CPU cores × 2)
	_defaultMaxIdleConns    = 10
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
