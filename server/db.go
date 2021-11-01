package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"gitlab.com/mirukakoro/sisikyo/db"

	"github.com/jmoiron/sqlx"
)

var driverName string
var dataSourceName string
var dbPing time.Duration

func init() {
	flag.StringVar(&driverName, "db-driver", "", "database: driver to use")
	flag.StringVar(&dataSourceName, "db-source", "", "database: data source")
	flag.DurationVar(&dbPing, "db-ping", 1*time.Second, "database: ping timeout")
}

func setupDb() (*sqlx.DB, error) {
	conn, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("conn: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbPing)
	defer cancel()
	err = conn.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	log.Print("db ping: ok")
	_, err = conn.Exec(db.Schema)
	if err != nil {
		return nil, fmt.Errorf("db schema: %w", err)
	}
	log.Print("db schema: ok")
	log.Printf("db: conn'd (driver: %s)", driverName)
	return conn, nil
}
