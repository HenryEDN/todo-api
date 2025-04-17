package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDB struct{
	DBPool *sql.DB
}

type DatabaseReader interface{
	CreateUser(user User) (*CreateUserResponseDTO, error)
	GetUsers() ([]*User, error)
}

const (
	GET_USERS = "SELECT * FROM users"
	GET_USER_BY_ID = "SELECT * FROM users WHERE userID = $1"
	CREATE_USER = "INSERT INTO users (userID, email, username, password, creation_time) VALUES ($1, $2, $3, $4, $5)"
)

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

func (db PostgresDB) CreateUser(user User) (*CreateUserResponseDTO, error){
	query := CREATE_USER

	_, err := db.DBPool.Exec(query, user.UserID, user.Email, user.Username, user.Password, user.Creation_time)
	if err != nil{
		return nil, err
	}

	return &CreateUserResponseDTO{Username: user.Username, Status: "Created"}, nil
}

func (db PostgresDB) GetUsers() ([]*User, error){
	query := GET_USERS

	users := []*User{}

	rows, err := db.DBPool.Query(query)
	if err != nil{
		return nil, err
	}

	for rows.Next(){
		user, err := scanIntoUser(rows)
		if err != nil{
			return nil, err
		}

		users = append(users, user)
	}

	fmt.Println(users)

	return users, nil
}


func scanIntoUser(rows *sql.Rows) (*User, error){
	user := new(User)
	err := rows.Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Creation_time)
	
	return user, err
}