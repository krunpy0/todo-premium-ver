package task

import (
	"time"
)

type Task struct {
	ID          string `json:"id"`
	UserID      string	`json:"user_id"`
	Title       string		`json:"title"`
	Description string		`json:"description"`
	Date        time.Time  `json:"date"`
	Difficulty  string     `json:"difficulty"`// easy / medium / hard — влияет на награду
	Status      string `json:"status"`// pending / done / failed
	CompletedAt *time.Time	`json:"completed_at"`
}
