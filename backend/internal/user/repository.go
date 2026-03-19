package user

import (
	"github.com/krunpy0/todo-premium-ver/db"
)

func QueryUserBool(username string) (bool, error) {
	var exists bool
	err := db.DB.QueryRow(`SELECT 1 WHERE EXISTS (SELECT 1 FROM users WHERE username =$1)`, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CreateUser(username string, hashedPassword string ) (int, error) {
	query := `
	INSERT INTO users (username, password)
	VALUES ($1, $2) RETURNING id
	`
	var id int
	err := db.DB.QueryRow(query, username, hashedPassword).Scan(id); if err != nil {
		return 0, err
	}
	return id, nil
}
