package achievement

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"context"
	"log"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type AchievementService struct {
	Exec boil.ContextExecutor
}

func (u *AchievementService) AddAchievementToUser(
	ctx context.Context,
	user_id string,
	achievements []*model.AchievementInput,
) error {

	db_user_achievement := &db.UserAchievement{}

	for _, achievement := range achievements {

		db_user_achievement = &db.UserAchievement{
			UserID:        user_id,
			AchievementID: achievement.ID,
		}
		err := db_user_achievement.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			log.Println("db_user_achievement.Insert error:", err)
			return err
		}
	}

	return nil
}

func (u *AchievementService) GetAchievementsByUserId(
	ctx context.Context,
	user_id string,
) ([]string, error) {
	result := []string{}
	
	db_user_achievements, err := db.UserAchievements(db.UserAchievementWhere.UserID.EQ(user_id)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_user_achievements.Select error:", err)
		return nil, err
	}

	for _, db_user_achievement := range db_user_achievements {
		result = append(result, db_user_achievement.AchievementID)
	}

	return result, nil
}
