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
WITH user_tz AS (
    SELECT timezone AS tz
    FROM users
    WHERE id = $1
),
today_local AS (
    SELECT (NOW() AT TIME ZONE tz)::date AS today
    FROM user_tz
)
UPDATE streaks s
SET
    current_streak = CASE
        WHEN s.last_active_date = tl.today THEN s.current_streak
        WHEN s.last_active_date = tl.today - 1 THEN s.current_streak + 1
        ELSE 1
    END,
    longest_streak = GREATEST(
        s.longest_streak,
        CASE
            WHEN s.last_active_date = tl.today THEN s.current_streak
            WHEN s.last_active_date = tl.today - 1 THEN s.current_streak + 1
            ELSE 1
        END
    ),
    last_active_date = tl.today,
    updated_at = NOW()
FROM today_local tl
WHERE s.user_id = $1
  AND (
    s.last_active_date IS NULL
    OR s.last_active_date <> tl.today
  );
	`

	_, err := db.DB.Exec(query, userID)
	return err
}

func RollbackStreakOnCancel(userID string) error {
	query := `
WITH user_tz AS (
    SELECT timezone AS tz
    FROM users
    WHERE id = $1
),
today_local AS (
    SELECT (NOW() AT TIME ZONE tz)::date AS today
    FROM user_tz
),
today_completions AS (
    SELECT COUNT(*) AS cnt
    FROM tasks t, user_tz u
    WHERE t.user_id = $1
      AND t.status = 'done'
      AND (t.completed_at AT TIME ZONE u.tz)::date = (NOW() AT TIME ZONE u.tz)::date
),
done_days AS (
    SELECT DISTINCT (t.completed_at AT TIME ZONE u.tz)::date AS done_day
    FROM tasks t, user_tz u
    WHERE t.user_id = $1
      AND t.status = 'done'
      AND t.completed_at IS NOT NULL
),
grouped_days AS (
    SELECT
        done_day,
        done_day - (ROW_NUMBER() OVER (ORDER BY done_day))::int AS grp
    FROM done_days
),
longest_calc AS (
    SELECT COALESCE(MAX(streak_len), 0) AS longest
    FROM (
        SELECT COUNT(*) AS streak_len
        FROM grouped_days
        GROUP BY grp
    ) s
)
UPDATE streaks s
SET
    current_streak = CASE
        WHEN tc.cnt = 0 THEN GREATEST(s.current_streak - 1, 0)
        ELSE s.current_streak
    END,
    last_active_date = CASE
        WHEN tc.cnt = 0 AND s.current_streak > 1 THEN tl.today - 1
        WHEN tc.cnt = 0 THEN NULL
        ELSE s.last_active_date
    END,
    longest_streak = CASE
        WHEN tc.cnt = 0 THEN GREATEST(
            lc.longest,
            CASE WHEN s.current_streak > 1 THEN s.current_streak - 1 ELSE 0 END
        )
        ELSE s.longest_streak
    END,
    updated_at = NOW()
FROM today_completions tc, today_local tl, longest_calc lc
WHERE s.user_id = $1
  AND s.last_active_date = tl.today;
	`

	_, err := db.DB.Exec(query, userID)
	return err
}
