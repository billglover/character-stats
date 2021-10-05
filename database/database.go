package database

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DSN string
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", cfg.DSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	// First check we can ping the database.
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity. Running this query forces a
	// round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

func NamedExecContext(ctx context.Context, log *log.Logger, db *sqlx.DB, query string, data interface{}) error {
	if _, err := db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	return nil
}
