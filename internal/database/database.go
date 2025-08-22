// Package database provides database connection management for the Mockbuster API.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq" //nolint:goimports //Recommended way to use the library
	"github.com/rxbenefits/go-hw/internal/util"
)

// DB holds the database connection.
type DB struct {
	*sql.DB
}

type dbOpts struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

type dbOptsFunc func(dbOpts) dbOpts

func defaultDBOpts() dbOpts {
	return dbOpts{
		host:     util.GetEnv("DB_HOST", "localhost"),
		port:     util.GetEnv("DB_PORT", "5555"),
		user:     util.GetEnv("DB_USER", "postgres"),
		password: util.GetEnv("DB_PASSWORD", "postgres"),
		dbname:   util.GetEnv("DB_NAME", "dvdrental"),
	}
}

// WithDBHost sets the database host.
func WithDBHost(host string) func(dbOpts) dbOpts {
	return func(opts dbOpts) dbOpts {
		opts.host = host
		return opts
	}
}

// WithDBPort sets the database port.
func WithDBPort(port string) func(dbOpts) dbOpts {
	return func(opts dbOpts) dbOpts {
		opts.port = port
		return opts
	}
}

// WithDBUser sets the database user.
func WithDBUser(user string) func(dbOpts) dbOpts {
	return func(opts dbOpts) dbOpts {
		opts.user = user
		return opts
	}
}

// WithDBPassword sets the database password.
func WithDBPassword(password string) func(dbOpts) dbOpts {
	return func(opts dbOpts) dbOpts {
		opts.password = password
		return opts
	}
}

// WithDBName sets the database name.
func WithDBName(dbname string) func(dbOpts) dbOpts {
	return func(opts dbOpts) dbOpts {
		opts.dbname = dbname
		return opts
	}
}

// InitDB initializes a new database connection with the given options.
func InitDB(opts ...dbOptsFunc) (*DB, error) {
	dbOptions := defaultDBOpts()

	for _, opt := range opts {
		opt(dbOptions)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbOptions.host, dbOptions.port, dbOptions.user, dbOptions.password, dbOptions.dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	slog.Info("Successfully connected to database")
	return &DB{db}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.DB.Close()
}
