package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDB struct{
	DBPool *sql.DB
}

func CreateDB(config Config) (*PostgresDB, error){
	connStr := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=%v", config.Username, config.DB_name, config.Password, config.SSLmode)

	dbPool, err := sql.Open("postgres", connStr)
	if err != nil{
		return nil, err
	}

	if err := dbPool.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{DBPool: dbPool}, nil
}