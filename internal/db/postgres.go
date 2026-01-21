package db

import (
	"database/sql"
	"log"
	"log/slog"
)

func New(dsn string, logger *slog.Logger) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to open database", "err", err)
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		logger.Error("failed to ping database", "err", err)
		log.Fatal(err)
	}

	logger.Info("Connected to DB")
	return db

}
