package main

import "time"

type User struct{
	UserID uint64 `json:"userID" db:"userID"`
	Email string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Creation_time time.Time `json:"creation_time" db:"creation_time"`
}

type CreateUserDTO struct{
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponseDTO struct{
	Username string `json:"username"`
	Status string `json:"status"`
}