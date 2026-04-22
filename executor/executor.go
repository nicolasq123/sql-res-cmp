package executor

import (
	"fmt"
	"net/url"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDB(dsn string) (*sqlx.DB, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	switch u.Scheme {
	case "mysql":
		db, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("connect: %w", err)
		}
		return db, nil
	case "clickhouse":
		db, err := sqlx.Connect("clickhouse", dsn)
		if err != nil {
			return nil, fmt.Errorf("connect: %w", err)
		}
		return db, nil
	}

	return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
}
