package user

import (
	"github.com/krunpy0/todo-premium-ver/db"
)

func QueryUserBool(username string) (bool, error) {
	var exists bool
	err := db.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE username=$1)`, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func QueryUser(username string) (User, error) {	
	var user = User{}
	err := db.DB.QueryRow(`SELECT id, username, password FROM users WHERE username=$1`, username).Scan(&user.ID,&user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CreateUser(username string, hashedPassword string ) (User, error) {
	var user = User{}
	query := `
	INSERT INTO users (username, password)
	VALUES ($1, $2) RETURNING id, username, password
	`
	err := db.DB.QueryRow(query, username, hashedPassword).Scan(&user.ID, &user.Username, &user.Password); if err != nil {
		return User{}, err
	}
	return user, nil
}

func UpdateUserXP(userID string, xpAmount int) (int, error) {
	var xp int
	if err := db.DB.QueryRow(`
		UPDATE users
		SET xp = xp + $1, updated_at = NOW()
		WHERE id = $2 RETURNING xp;`, xpAmount, userID).Scan(&xp); err != nil {
		return 0, err
	}
	return xp, nil
}

func UpdateUserCoins(userID string, coinsAmount int) (int, error) {
	var coins int
	if err := db.DB.QueryRow(`
		UPDATE users
		SET coins = coins + $1, updated_at = NOW()
		WHERE id = $2 RETURNING coins;`, coinsAmount, userID).Scan(&coins); err != nil {
		return 0, err
	}
	return coins, nil
}

