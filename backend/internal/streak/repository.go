package streak

import (
	"github.com/krunpy0/todo-premium-ver/db"
)

func QueryStreak(userID string) (Streak, error) {
	var streak = Streak{}

	err := db.DB.QueryRow(`SELECT user_id, current_streak, longest_streak, last_active_date  FROM streaks WHERE user_id=$1`, userID).Scan(&streak.UserID, &streak.CurrentStreak, &streak.LongestStreak, &streak.LastActiveDate)
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
	return err
}

func UpdateStreak(userID string) error {
	query := `
	UPDATE streaks
SET
    current_streak = CASE
        WHEN last_active_date = CURRENT_DATE AT TIME ZONE (SELECT timezone FROM users WHERE id = user_id)
            THEN current_streak  -- уже обновляли сегодня, ничего не меняем

        WHEN last_active_date = (CURRENT_DATE AT TIME ZONE (SELECT timezone FROM users WHERE id = user_id)) - INTERVAL '1 day'
            THEN current_streak + 1  -- вчера было действие — продолжаем стрик

        ELSE 1  -- первый день или перерыв — новый стрик с 1
    END,

    longest_streak = GREATEST(
        longest_streak,
        CASE
            WHEN last_active_date = (CURRENT_DATE AT TIME ZONE (SELECT timezone FROM users WHERE id = user_id)) - INTERVAL '1 day'
                THEN current_streak + 1
            ELSE 1
        END
    ),

    last_active_date = CURRENT_DATE AT TIME ZONE (SELECT timezone FROM users WHERE id = user_id),
    updated_at = NOW()

WHERE user_id = $1
  AND (
    last_active_date IS NULL
    OR last_active_date <> CURRENT_DATE AT TIME ZONE (SELECT timezone FROM users WHERE id = user_id)
  );
	`

	_, err := db.DB.Exec(query, userID)
	return err
}
