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
			Category: &model.AchievementCategory{
				ID:          db_achievement_category.ID,
				Name:        db_achievement_category.Name,
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
			ID:          db_achievement_category.ID,
			Name:        db_achievement_category.Name,
			DisplayName: db_achievement_category.DisplayName,
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
			ID:          db_achievement_category.ID,
			Name:        db_achievement_category.Name,
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

func (u *AchievementService) AddChillAchievement(
	ctx context.Context,
	user_id string,
	chill_id string,
	achievement_id []string,
) error {

	have_achievements, err := db.UserAchievements(db.UserAchievementWhere.UserID.EQ(user_id)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_user_achievements.Select error:", err)
		return err
	}

	is_have_list := map[string]bool{}

	for _, id := range achievement_id {
		is_have_list[id] = false
		for _, have_achievement := range have_achievements {
			if have_achievement.AchievementID == id {
				is_have_list[id] = true
				break
			}
		}
	}

	for _, id := range achievement_id {
		if !is_have_list[id] {
			db_chill_achievement := &db.ChillAchievement{
				ChillID:       chill_id,
				AchievementID: id,
			}

			err := db_chill_achievement.Insert(ctx, u.Exec, boil.Infer())
			if err != nil {
				log.Println("db_chill_achievement.Insert error:", err)
				return err
			}

			db_user_achievement := &db.UserAchievement{
				UserID:        user_id,
				AchievementID: id,
			}

			err = db_user_achievement.Insert(ctx, u.Exec, boil.Infer())
			if err != nil {
				log.Println("db_user_achievement.Insert error:", err)
				return err
			}
		}
	}

	return nil
}

func (u *AchievementService) GetAvatarByUser(
	ctx context.Context,
	user *model.User,
) (*model.Achievement, error) {
	if user.Avatar == nil {
		return &model.Achievement{}, nil
	}

	db_achievement, err := db.Achievements(db.AchievementWhere.ID.EQ(user.Avatar.ID)).One(ctx, u.Exec)
	if err != nil {
		log.Println("db_achievement.Select error:", err)
		return &model.Achievement{}, err
	}

	return &model.Achievement{
		ID:          db_achievement.ID,
		Name:        db_achievement.Name,
		DisplayName: db_achievement.DisplayName,
		Description: db_achievement.Description,
	}, nil
}

func (u *AchievementService) GetAchievementCategoryByAchievement(
	ctx context.Context,
	achievement *model.Achievement,
) (*model.AchievementCategory, error) {
	if achievement == nil {
		return &model.AchievementCategory{}, nil
	}

	db_achievement_category, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(achievement.Category.ID)).One(ctx, u.Exec)
	if err != nil {
		log.Println("db_achievement_category.Select error:", err)
		return &model.AchievementCategory{}, err
	}

	return &model.AchievementCategory{
		ID:          db_achievement_category.ID,
		Name:        db_achievement_category.Name,
		DisplayName: db_achievement_category.DisplayName,
	}, nil
}

func (u *AchievementService) GetAchievementsByAchievementCategory(
	ctx context.Context,
	achievement_category *model.AchievementCategory,
) ([]*model.Achievement, error) {
	if achievement_category == nil {
		return []*model.Achievement{}, nil
	}

	result := []*model.Achievement{}

	db_achievements, err := db.Achievements(db.AchievementWhere.CategoryID.EQ(achievement_category.ID)).All(ctx, u.Exec)
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

		result = append(result, achievement)
	}

	return result, nil
}
