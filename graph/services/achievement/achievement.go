package achievement

import (
	"chilly_daze_gateway/graph/db"
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type AchievementService struct {
	Exec boil.ContextExecutor
}

func (u *AchievementService) AddAchievementToUser(
	ctx context.Context,
	user_id string,
	achievementID string,
) error {
	db_user_achievement := &db.UserAchievement{
		UserID:        user_id,
		AchievementID: achievementID,
	}

	err := db_user_achievement.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}
