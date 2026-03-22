package streak

import (
	"time"
)

type Streak struct {
	UserID string `json:"user_id"`
	CurrentStreak int `json:"current_streak"`
	LongestStreak int `json:"longest_streak"`
	LastActiveDate 		*time.Time `json:"last_active_date"`
}