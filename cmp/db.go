package cmp

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// parseMySQLURLToDSN converts mysql:// URL format to Go MySQL driver DSN format
func parseMySQLURL(dsn string) string {
	if !strings.HasPrefix(dsn, "mysql://") {
		return dsn
	}

	if !strings.Contains(dsn, "@tcp(") {
		return dsn
	}

	dsn = strings.ReplaceAll(dsn, "@tcp(", "@")
	dsn = strings.ReplaceAll(dsn, ")", "")
	return dsn
}

func NewDB(dsn string) (*sqlx.DB, error) {
	// Special handling for mysql with tcp(...) format
	dsn = parseMySQLURL(dsn)

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	var driver, realDSN string
	switch u.Scheme {
	case "mysql":
		host := u.Hostname()
		port := u.Port()
		if port == "" {
			port = "3306"
		}
		user := u.User.Username()
		pass, _ := u.User.Password()
		dbname := u.Path[1:]
		driver = "mysql"
		realDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)
	case "clickhouse":
		host := u.Hostname()
		port := u.Port()
		if port == "" {
			port = "9000"
		}
		q := u.Query()
		user := q.Get("username")
		pass := q.Get("password")
		dbname := u.Path[1:]
		driver = "clickhouse"
		realDSN = fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s", host, port, user, pass, dbname)
	case "postgres":
		host := u.Hostname()
		port := u.Port()
		if port == "" {
			port = "5432"
		}
		user := u.User.Username()
		pass, _ := u.User.Password()
		dbname := u.Path[1:]
		driver = "postgres"
		realDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	db, err := sqlx.Connect(driver, realDSN)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(0)
	return db, nil
}

func Query(ctx context.Context, db *sqlx.DB, q string, args ...any) ([]string, [][]string, error) {
	rows, err := db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var result [][]string
	for rows.Next() {
		dests := make([]any, len(cols))
		byteSlices := make([][]byte, len(cols))
		for i := range dests {
			dests[i] = &byteSlices[i]
		}
		if err := rows.Scan(dests...); err != nil {
			return nil, nil, err
		}
		row := make([]string, len(cols))
		for i, b := range byteSlices {
			if b == nil {
				row[i] = "NULL"
			} else {
				row[i] = string(b)
			}
		}
		result = append(result, row)
	}
	return cols, result, rows.Err()
}
