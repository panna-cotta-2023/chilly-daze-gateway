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

		db_achievement_category, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(db_achievement.CategoryID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement_category.Select error:", err)
			return nil, err
		}

		result = append(result, &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			Description: db_achievement.Description,
			DisplayName: db_achievement_category.DisplayName,
			Category: 	&model.AchievementCategory{
				ID: db_achievement_category.ID,
				Name: db_achievement_category.Name,
				DisplayName: db_achievement_category.DisplayName,
			},
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
		achievement := &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			DisplayName: db_achievement.DisplayName,
			Description: db_achievement.Description,
		}

		db_achievement_category, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(db_achievement.CategoryID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement_category.Select error:", err)
			return nil, err
		}

		achievement.Category = &model.AchievementCategory{
			ID: db_achievement_category.ID,
			Name: db_achievement_category.Name,
			DisplayName: db_achievement_category.DisplayName,
		}

		db_achievements_achievement_category, err := db.Achievements(db.AchievementWhere.CategoryID.EQ(db_achievement_category.ID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievements_category.Select error:", err)
			return nil, err
		}

		for _, db_achievement_achievement_category := range db_achievements_achievement_category {
			achievement.Category.Achievements = append(achievement.Category.Achievements, &model.Achievement{
				ID:          db_achievement_achievement_category.ID,
				Name:        db_achievement_achievement_category.Name,
				DisplayName: db_achievement_achievement_category.DisplayName,
				Description: db_achievement_achievement_category.Description,
			})
		}

		result = append(result, achievement)
	}

	return result, nil
}

func (u *AchievementService) GetAchievementCategories(
	ctx context.Context,
) ([]*model.AchievementCategory, error) {

	result := []*model.AchievementCategory{}

	db_achievement_categories, err := db.AchievementCategories().All(ctx, u.Exec)
	if err != nil {
		log.Println("db_achievement_categories.Select error:", err)
		return nil, err
	}

	for _, db_achievement_category := range db_achievement_categories {
		achievement_category := &model.AchievementCategory{
			ID: db_achievement_category.ID,
			Name: db_achievement_category.Name,
			DisplayName: db_achievement_category.DisplayName,
		}

		db_achievements, err := db.Achievements(db.AchievementWhere.CategoryID.EQ(db_achievement_category.ID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievements.Select error:", err)
			return nil, err
		}

		for _, db_achievement := range db_achievements {
			achievement_category.Achievements = append(achievement_category.Achievements, &model.Achievement{
				ID:          db_achievement.ID,
				Name:        db_achievement.Name,
				DisplayName: db_achievement.DisplayName,
				Description: db_achievement.Description,
			})
		}

		result = append(result, achievement_category)
	}
	
	return result, nil
}

// func (u *AchievementService) GetAchievementsByCategoryId(
// 	ctx context.Context,
// 	category_id string,
// ) ([]*model.Achievement, error) {
// 	result := []*model.Achievement{}

// 	db_achievements, err := db.Achievements(db.AchievementWhere.CategoryID.EQ(category_id)).All(ctx, u.Exec)
// 	if err != nil {
// 		log.Println("db_achievements.Select error:", err)
// 		return nil, err
// 	}

// 	for _, db_achievement := range db_achievements {
// 		result = append(result, &model.Achievement{
// 			ID:          db_achievement.ID,
// 			Name:        db_achievement.Name,
// 			DisplayName: db_achievement.DisplayName,
// 			Description: db_achievement.Description,
// 		})
// 	}

// 	return result, nil
// }

// func (u *AchievementService) GetAchievementCategoryByAchievementId(
// 	ctx context.Context,
// 	achievement_id string,
// ) (*model.AchievementCategory, error) {
// 	result := &model.AchievementCategory{}

// 	db_achievement, err := db.Achievements(db.AchievementWhere.ID.EQ(achievement_id)).One(ctx, u.Exec)
// 	if err != nil {
// 		log.Println("db_achievement.Select error:", err)
// 		return nil, err
// 	}

// 	db_achievement_category, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(db_achievement.CategoryID)).One(ctx, u.Exec)
// 	if err != nil {
// 		log.Println("db_achievement_category.Select error:", err)
// 		return nil, err
// 	}

// 	result.ID = db_achievement_category.ID
// 	result.Name = db_achievement_category.Name
// 	result.DisplayName = db_achievement_category.DisplayName

// 	return result, nil
// }