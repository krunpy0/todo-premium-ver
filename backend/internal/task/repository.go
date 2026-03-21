package task

import (
	"time"

	db "github.com/krunpy0/todo-premium-ver/db"
)

func GetUserTasks(userID string) ([]Task, error) {
	var tasks = []Task{}
	rows, err := db.DB.Query(`SELECT * FROM tasks WHERE user_id = $1`, userID)
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var t = Task{}
		err := rows.Scan(
			&t.ID, 
			&t.UserID, 
			&t.Title, 
			&t.Date, 
			&t.Difficulty, 
			&t.Status, 
			&t.CompletedAt,
			&t.Due,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func CreateTask(userID string, title string, difficulty string, due *time.Time) (string, error) {
	var id string

	if err := db.DB.QueryRow(`INSERT INTO tasks (user_id, title, date, difficulty, due)
	VALUES ($1, $2, NOW(), $3, $4) RETURNING id`, userID, title, difficulty, due).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}