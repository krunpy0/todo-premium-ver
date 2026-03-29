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

func GetTaskById(taskID string, userID string) (Task, error) {
	var task = Task{}
	err := db.DB.QueryRow(`SELECT * FROM tasks WHERE id=$1 AND user_id=$2`, taskID, userID).Scan(&task.ID, &task.UserID, &task.Title, &task.Date, &task.Difficulty, &task.Status, &task.CompletedAt, &task.Due)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func CreateTask(userID string, title string, difficulty string, due *time.Time) (string, error) {
	var id string

	if err := db.DB.QueryRow(`INSERT INTO tasks (user_id, title, date, difficulty, due)
	VALUES ($1, $2, NOW(), $3, $4) RETURNING id`, userID, title, difficulty, due).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func CompleteTask(userID string, taskID string) (Task, error) {
	var task = Task{}

	if err := db.DB.QueryRow(`UPDATE tasks
		SET status = 'done', completed_at = NOW()
		WHERE id = $1 AND user_id = $2 AND status != 'done'
		RETURNING *`, taskID, userID).Scan(&task.ID, &task.UserID, &task.Title, &task.Date, &task.Difficulty, &task.Status, &task.CompletedAt, &task.Due); err != nil {
		return Task{}, err
	}
	return task, nil
}

func FailTask(userID string, taskID string) (Task, error) {
	var task = Task{}

	if err := db.DB.QueryRow(`UPDATE tasks
		SET status = 'failed', completed_at = NOW()
		WHERE id = $1 AND user_id = $2 AND status != 'failed' AND status != 'done' 
		RETURNING *`, taskID, userID).Scan(&task.ID, &task.UserID, &task.Title, &task.Date, &task.Difficulty, &task.Status, &task.CompletedAt, &task.Due); err != nil {
		return Task{}, err
	}
	return task, nil
}

func CancelTask(userID string, taskID string) (Task, error) {
	var task = Task{}

	if err := db.DB.QueryRow(`UPDATE tasks
		SET status = 'pending', completed_at = NULL
		WHERE id = $1 AND user_id = $2 AND (status = 'done' OR status = 'failed')
		RETURNING *`, taskID, userID).Scan(&task.ID, &task.UserID, &task.Title, &task.Date, &task.Difficulty, &task.Status, &task.CompletedAt, &task.Due); err != nil {
		return Task{}, err
	}
	return task, nil
}
