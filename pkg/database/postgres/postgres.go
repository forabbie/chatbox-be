package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB   *sql.DB
	Conn *sql.Conn
}

func Open(user, pass, host, port, name, sslmode string, ctx context.Context) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user,
		pass,
		host,
		port,
		name,
		sslmode,
	)

	db, _ := sql.Open("postgres", dsn)

	conn, err := db.Conn(ctx)

	pg := new(PostgresDB)

	pg.DB = db

	pg.Conn = conn

	return pg, err
}

func (p *PostgresDB) Close() error {
	_ = p.Conn.Close()
	return p.DB.Close()
}

func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.Conn.PingContext(ctx)
}
