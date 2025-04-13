package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLd struct {
	DB *sql.DB
	Conn *sql.Conn
}

func Open(user, pass, host, port, name string, ctx context.Context) (*MySQLd, error) {
	dsn := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=%s&collation=%s&interpolateParams=%t&loc=%s&multiStatements=%t&parseTime=%t&sql_mode=%s&time_zone=%s",
		user,
		pass,
		"tcp",
		host,
		port,
		name,
		"utf8mb4,utf8",
		"utf8mb4_general_ci",
		false,
		time.UTC,
		false,
		true,
		"%27%27", // ''
		"%27%2B00%3A00%27", // '+00:00'
	)

	db, _ := sql.Open("mysql", dsn)

	conn, err := db.Conn(ctx)

	mysqld := new(MySQLd)

	mysqld.DB = db

	mysqld.Conn = conn

	return mysqld, err
}

func (m *MySQLd) Close() error {
	_ = m.Conn.Close()

	err := m.DB.Close()

	return err
}

func (m *MySQLd) Ping(ctx context.Context) error {
	return m.Conn.PingContext(ctx)
}
