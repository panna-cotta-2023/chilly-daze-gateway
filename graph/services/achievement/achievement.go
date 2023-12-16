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
	userId string,
) ([]*model.Achievement, error) {
	result := []*model.Achievement{}

	dbUserAchievements, err := db.UserAchievements(db.UserAchievementWhere.UserID.EQ(userId)).All(ctx, u.Exec)
	if err != nil {
		log.Println("dbUserAchievements.Select error:", err)
		return nil, err
	}

	for _, dbUserAchievement := range dbUserAchievements {
		dbAchievements, err := db.Achievements(db.AchievementWhere.ID.EQ(dbUserAchievement.AchievementID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("dbAchievement.Select error:", err)
			return nil, err
		}

		for _, dbAchievement := range dbAchievements {
			dbAchievementCategory, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(dbAchievement.CategoryID)).One(ctx, u.Exec)
			if err != nil {
				log.Println("dbAchievementCategory.Select error:", err)
				return nil, err
			}

			result = append(result, &model.Achievement{
				ID:          dbAchievement.ID,
				Name:        dbAchievement.Name,
				Description: dbAchievement.Description,
				DisplayName: dbAchievementCategory.DisplayName,
				Category: &model.AchievementCategory{
					ID:          dbAchievementCategory.ID,
					Name:        dbAchievementCategory.Name,
					DisplayName: dbAchievementCategory.DisplayName,
				},
			})
		}

	}

	return result, nil
}

func (u *AchievementService) GetAchievements(
	ctx context.Context,
) ([]*model.Achievement, error) {

	result := []*model.Achievement{}

	dbAchievements, err := db.Achievements().All(ctx, u.Exec)
	if err != nil {
		log.Println("dbAchievements.Select error:", err)
		return nil, err
	}

	for _, dbAchievement := range dbAchievements {
		achievement := &model.Achievement{
			ID:          dbAchievement.ID,
			Name:        dbAchievement.Name,
			DisplayName: dbAchievement.DisplayName,
			Description: dbAchievement.Description,
		}

		dbAchievementCategory, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(dbAchievement.CategoryID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("dbAchievementCategory.Select error:", err)
			return nil, err
		}

		achievement.Category = &model.AchievementCategory{
			ID:          dbAchievementCategory.ID,
			Name:        dbAchievementCategory.Name,
			DisplayName: dbAchievementCategory.DisplayName,
		}

		result = append(result, achievement)
	}

	return result, nil
}

func (u *AchievementService) GetAchievementCategories(
	ctx context.Context,
) ([]*model.AchievementCategory, error) {

	result := []*model.AchievementCategory{}

	dbAchievementCategories, err := db.AchievementCategories().All(ctx, u.Exec)
	if err != nil {
		log.Println("dbAchievementCategories.Select error:", err)
		return nil, err
	}

	for _, dbAchievementCategory := range dbAchievementCategories {
		achievementCategory := &model.AchievementCategory{
			ID:          dbAchievementCategory.ID,
			Name:        dbAchievementCategory.Name,
			DisplayName: dbAchievementCategory.DisplayName,
		}

		dbAchievements, err := db.Achievements(db.AchievementWhere.CategoryID.EQ(dbAchievementCategory.ID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("dbAchievements.Select error:", err)
			return nil, err
		}

		for _, dbAchievement := range dbAchievements {
			achievementCategory.Achievements = append(achievementCategory.Achievements, &model.Achievement{
				ID:          dbAchievement.ID,
				Name:        dbAchievement.Name,
				DisplayName: dbAchievement.DisplayName,
				Description: dbAchievement.Description,
			})
		}

		result = append(result, achievementCategory)
	}

	return result, nil
}

func (u *AchievementService) GetAvatarByUser(
	ctx context.Context,
	user *model.User,
) (*model.Achievement, error) {
	if user.Avatar == nil || user.Avatar.ID == "" {
		return &model.Achievement{}, nil
	}

	dbAchievements, err := db.Achievements(db.AchievementWhere.ID.EQ(user.Avatar.ID)).All(ctx, u.Exec)
	if err != nil {
		log.Println("dbAchievement.Select error:", err)
		return &model.Achievement{}, err
	}

	for _, dbAchievement := range dbAchievements {
		return &model.Achievement{
			ID:          dbAchievement.ID,
			Name:        dbAchievement.Name,
			Description: dbAchievement.Description,
			DisplayName: dbAchievement.DisplayName,
			Category: &model.AchievementCategory{
				ID:          dbAchievement.CategoryID,
			},
		}, nil
	}
	return &model.Achievement{}, nil
}

func (u *AchievementService) GetAchievementCategoryByAchievement(
	ctx context.Context,
	achievement *model.Achievement,
) (*model.AchievementCategory, error) {
	if achievement == nil || achievement.Category == nil {
		return &model.AchievementCategory{}, nil
	}

	dbAchievementCategory, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(achievement.Category.ID)).One(ctx, u.Exec)
	if err != nil {
		log.Println("dbAchievementCategory.Select error:", err)
		return &model.AchievementCategory{}, err
	}

	return &model.AchievementCategory{
		ID:          dbAchievementCategory.ID,
		Name:        dbAchievementCategory.Name,
		DisplayName: dbAchievementCategory.DisplayName,
	}, nil
}

func (u *AchievementService) GetAchievementsByAchievementCategory(
	ctx context.Context,
	achievementCategory *model.AchievementCategory,
) ([]*model.Achievement, error) {
	if achievementCategory == nil {
		return []*model.Achievement{}, nil
	}

	result := []*model.Achievement{}

	dbAchievements, err := db.Achievements(db.AchievementWhere.CategoryID.EQ(achievementCategory.ID)).All(ctx, u.Exec)
	if err != nil {
		log.Println("dbAchievements.Select error:", err)
		return nil, err
	}

	for _, dbAchievement := range dbAchievements {
		achievement := &model.Achievement{
			ID:          dbAchievement.ID,
			Name:        dbAchievement.Name,
			DisplayName: dbAchievement.DisplayName,
			Description: dbAchievement.Description,
		}

		result = append(result, achievement)
	}

	return result, nil
}

func (u *AchievementService) GetNewAchievements(
	ctx context.Context,
	chill *model.Chill,
	userId string,
) ([]*model.Achievement, error) {
	// ToDo: check achievement
	achievementIds := []string{"423a969b-76bd-4848-88bf-9f6bf494fdc7"}

	result := []*model.Achievement{}

	userAchievements, err := db.UserAchievements(
		db.UserAchievementWhere.UserID.EQ(userId),
	).All(ctx, u.Exec)
	if err != nil {
		log.Println("db.UserAchievements error:", err)
		return nil, err
	}

	newAchievementIds := []string{}

	for _, achievementId := range achievementIds {
		flag := false
		for _, userAchievement := range userAchievements {
			if userAchievement.AchievementID == achievementId {
				flag = true
			}
		}
		if !flag {
			newAchievementIds = append(newAchievementIds, achievementId)
		}
	}

	dbNewAchievements, err := db.Achievements(
		db.AchievementWhere.ID.IN(newAchievementIds),
	).All(ctx, u.Exec)
	if err != nil {
		log.Println("db.Achievements error:", err)
		return nil, err
	}

	for _, dbNewAchievement := range dbNewAchievements {

		dbChillAchievement := &db.ChillAchievement{
			ChillID:       chill.ID,
			AchievementID: dbNewAchievement.ID,
		}

		if err = dbChillAchievement.Insert(ctx, u.Exec, boil.Infer()); err != nil {
			log.Println("dbChillAchievement.Insert error:", err)
			return nil, err
		}

		dbUserAchievement := &db.UserAchievement{
			UserID:        userId,
			AchievementID: dbNewAchievement.ID,
		}

		if err = dbUserAchievement.Insert(ctx, u.Exec, boil.Infer()); err != nil {
			log.Println("db_user_achievement.Insert error:", err)
			return nil, err
		}
	}

	dbChillAchievemnts, err := db.ChillAchievements(db.ChillAchievementWhere.ChillID.EQ(chill.ID)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db.ChillAchievements error:", err)
		return nil, err
	}

	for _, dbChillAchievement := range dbChillAchievemnts {
		dbAchievements, err := db.Achievements(db.AchievementWhere.ID.EQ(dbChillAchievement.AchievementID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db.Achievements error:", err)
			return nil, err
		}

		for _, dbAchievement := range dbAchievements {
			result = append(result, &model.Achievement{
				ID:          dbAchievement.ID,
				Name:        dbAchievement.Name,
				DisplayName: dbAchievement.DisplayName,
				Description: dbAchievement.Description,
				Category: &model.AchievementCategory{
					ID:          dbAchievement.CategoryID,
				},
			})
		}
	}

	return result, nil
}
