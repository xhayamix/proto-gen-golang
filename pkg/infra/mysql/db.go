package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
)

const (
	defaultMaxIdleConns = 2
	defaultMaxOpenConns = 100

	// PlaceholderLimit プリペアドステートメントの上限
	// https://blog.pinkumohikan.com/entry/workaround-for-mysql-too-many-placeholders-error
	PlaceholderLimit = 65535
)

type Config struct {
	Addr     string
	User     string
	Password string
	DB       string
}

func New(c *Config) (*sql.DB, error) {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local", c.User, c.Password, c.Addr, c.DB)
	db, err := sql.Open("transaction", dataSourceName)
	if err != nil {
		return nil, cerrors.Wrapf(err, cerrors.Internal, "MySQLに接続できませんでした。 dataSourceName = %s", dataSourceName)
	}
	db.SetMaxIdleConns(defaultMaxIdleConns)
	db.SetMaxOpenConns(defaultMaxOpenConns)
	db.SetConnMaxLifetime(defaultMaxOpenConns * time.Second)
	// TODO: クローズ処理実装

	return db, nil
}

type MysqlDB *sql.DB
