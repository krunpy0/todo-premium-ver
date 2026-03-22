package streak

import (
	"github.com/krunpy0/todo-premium-ver/db"
)

func QueryStreak(userID string) (Streak, error) {
	var streak = Streak{}

	err := db.DB.QueryRow(`SELECT user_id, current_streak, longset_streak, last_active_date  FROM streaks WHERE user_id=$1`, userID).Scan(&streak.UserID, &streak.CurrentStreak, &streak.LongestStreak, &streak.LastActiveDate)
	if err != nil {
		return Streak{}, err
	}
	return streak, nil
}

func CreateStreak(userID string) error {

	query := `
	INSERT INTO streaks (user_id)
	VALUES ($1)
	`
	_, err := db.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}