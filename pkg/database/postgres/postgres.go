package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	User, Pass, Host, Port, Name string
	SSLMode, TimeZone            string
}

type PostgresDB struct {
	DB   *sql.DB
	Conn *sql.Conn
}

func Open(ctx context.Context, config Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s TimeZone=%s",
		config.User,
		config.Pass,
		config.Host,
		config.Port,
		config.Name,
		config.SSLMode,
		config.TimeZone,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		DB:   db,
		Conn: conn,
	}, nil
}

func (p *PostgresDB) Close() error {
	_ = p.Conn.Close()
	return p.DB.Close()
}

func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.Conn.PingContext(ctx)
}
