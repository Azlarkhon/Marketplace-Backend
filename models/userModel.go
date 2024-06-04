package models

import (
	"time"
)

type User struct {
	ID         int       `json:"user_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Address    string    `json:"address"`
	Password   string    `json:"password"`
	Age        int       `json:"age"`
	Created_at time.Time `json:"created_at"`
}
