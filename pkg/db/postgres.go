package db

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"property-managment-service/internal/config"
)

func NewPsqlDB(cfg *config.Config) (*sqlx.DB, error) {
	const op = "db.NewPsqlDB"
	var sslMode string
	if cfg.Postgres.SslMode == true {
		sslMode = "require"
	} else {
		sslMode = "disable"
	}

	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.DbName,
		sslMode,
		cfg.Postgres.Password,
	)

	db, err := sqlx.Connect("pgx", dataSourceName)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
