package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Config :
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

const (
	postgres = "postgres"
)

// NewPostgresDB : создает подключение к базе данных
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s sslmode=%s",
		cfg.Host, cfg.Username, cfg.DBName, cfg.Port, cfg.Password, cfg.SSLMode)

	db, err := sqlx.Open(postgres, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
