package user

import (
	"time"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	XP        int `json:"xp"`
	Coins     int `json:"coins"`
	TimeZone 	string `json:"timezone"`
}