package datastore

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kadeallendev/bookstore/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(cfg config.Config) (*sql.DB, error) {
	options := fmt.Sprintf("user=%v password=%v host=%v port=%v database=%v sslmode=disable", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("pgx", options)
	if err != nil {
		errMsg := fmt.Sprintf("error opening DB connection: %v", err)
		return nil, errors.New(errMsg)
	}

	if err = db.Ping(); err != nil {
		errMsg := fmt.Sprintf("error pinging DB: %v", err)
		return nil, errors.New(errMsg)
	}

	return db, nil
}
