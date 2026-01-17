package db

import (
	"database/sql"
	"log"
)

func New(dsn string) *sql.DB {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")
	return db

}
