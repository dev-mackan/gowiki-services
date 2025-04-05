package db

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string //e.g. "15m"
}

func New(dbCfg *DbConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbCfg.Addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbCfg.MaxOpenConns)
	db.SetMaxIdleConns(dbCfg.MaxIdleConns)
	//TODO: get duration from env
	duration, err := time.ParseDuration(dbCfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
