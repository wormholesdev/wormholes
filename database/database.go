package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

const INITIAL_SCHEMA = `
-- enable uuid
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- users
CREATE TABLE IF NOT EXISTS wh_users (
  id uuid DEFAULT uuid_generate_v4(),
  email text NOT NULL UNIQUE,
  hashed_password varchar(255) NOT NULL,
  created_at timestamptz NOT NULL default now(),
  CONSTRAINT wh_users_pk PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS wh_links (
  id varchar(255),
  tag varchar(255),
  target varchar(255),
  created_at timestamptz NOT NULL default now(),
  CONSTRAINT wh_links_pk PRIMARY KEY (id)
);
`

// TODO: Add schema for timeseries events

// Postgres connection
type Postgres struct {
	Username string
	Password string
	Host     string
	Database string
	Port     int
}

func Default() Postgres {
	return Postgres{
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     5432,
		Database: "postgres",
	}
}

func (p *Postgres) Connect() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", p.Username, p.Password, p.Host, p.Port, p.Database),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	config.MaxConns = 200
	dbpool, err := pgxpool.ConnectConfig(
		context.Background(), config,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	_, err = dbpool.Exec(context.Background(), INITIAL_SCHEMA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create required tables: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
