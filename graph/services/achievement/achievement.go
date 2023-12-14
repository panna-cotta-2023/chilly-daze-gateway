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
	get_achievements []*model.AchievementInput,
	having_achievements []*model.Achievement,
) error {
	isHavingList := map[string]bool{}

	for _, get_achievement := range get_achievements {
		isHavingList[get_achievement.ID] = false

		for _, having_achievement := range having_achievements {

			if having_achievement.ID == get_achievement.ID {
				isHavingList[get_achievement.ID] = true
				break
			}
		}
	}

	for _, get_achievement := range get_achievements {
		if isHavingList[get_achievement.ID] {
			continue
		}

		db_user_achievement := &db.UserAchievement{
			UserID:       user_id,
			AchievementID: get_achievement.ID,
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
) ([]*model.Achievement, error) {
	result := []*model.Achievement{}

	db_user_achievements, err := db.UserAchievements(db.UserAchievementWhere.UserID.EQ(user_id)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_user_achievements.Select error:", err)
		return nil, err
	}

	for _, db_user_achievement := range db_user_achievements {
		db_achievement, err := db.Achievements(db.AchievementWhere.ID.EQ(db_user_achievement.AchievementID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement.Select error:", err)
			return nil, err
		}

		result = append(result, &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			Description: db_achievement.Description.String,
			// Category:    db_achievement.Category,
			// Image: 		 db_achievement.Image,
		})
	}

	return result, nil
}

func (u *AchievementService) GetAchievements(
	ctx context.Context,
) ([]*model.Achievement, error) {
	
	result := []*model.Achievement{}

	db_achievements, err := db.Achievements().All(ctx, u.Exec)
	if err != nil {
		log.Println("db_achievements.Select error:", err)
		return nil, err
	}

	for _, db_achievement := range db_achievements {
		result = append(result, &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			Description: db_achievement.Description.String,
			// Category:    db_achievement.Category,
			// Image: 		 db_achievement.Image,
		})
	}

	return result, nil
}