package main

import (
	"database/sql"
	"log"

	"github.com/pressly/goose"
)

func MigrateUp(db *sql.DB) error{
	log.Println("migrating...")
	return goose.Up(db, "./migrations")
}