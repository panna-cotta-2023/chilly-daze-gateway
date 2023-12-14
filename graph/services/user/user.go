package user

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"context"
	"log"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserService struct {
	Exec boil.ContextExecutor
}

func (u *UserService) CreateUser(
	ctx context.Context,
	input model.RegisterUserInput,
	uid string,
) (*model.User, error) {

	name := ""
	avatar := ""

	if input.Name != nil {
		name = *input.Name
	}

	if input.Avatar != nil {
		avatar = *input.Avatar
	}

	result := &model.User{}

	db_user := &db.User{
		ID:        uid,
		Name:      name,
		Avatar:    null.StringFromPtr(input.Avatar),
		CreatedAt: time.Now(),
	}

	result.ID = uid
	result.Name = name

	log.Println("avatar: ", avatar)

	db_achievements, err := db.Achievements(db.AchievementWhere.Name.EQ(avatar)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_achievement.Select error:", err)
	}

	for _, db_achievement := range db_achievements {
		result.Avatar = &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			Description: db_achievement.Description,
			DisplayName: db_achievement.DisplayName,
		}

		db_user.Avatar = null.StringFrom(db_achievement.ID)

		db_achievement_category, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(db_achievement.CategoryID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement_category.Select error:", err)

		}

		result.Avatar.Category = &model.AchievementCategory{
			ID:          db_achievement_category.ID,
			Name:        db_achievement_category.Name,
			DisplayName: db_achievement_category.DisplayName,
		}
	}

	db_user_chills, err := db.UserChills(db.UserChillWhere.UserID.EQ(uid)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_user_chills.Select error:", err)

	}

	for _, db_user_chill := range db_user_chills {
		db_chill, err := db.Chills(db.ChillWhere.ID.EQ(db_user_chill.ChillID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_chill.Select error:", err)

		}

		chill := &model.Chill{
			ID: db_chill.ID,
		}

		db_traces, err := db.TracePoints(db.TracePointWhere.ChillID.EQ(db_chill.ID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db_traces.Select error:", err)

		}

		for _, db_trace := range db_traces {
			trace := &model.TracePoint{
				ID:        db_trace.ID,
				Timestamp: db_trace.Timestamp.Format("2006-01-02T15:04:05+09:00"),
				Coordinate: &model.Coordinate{
					Latitude:  db_trace.Latitude,
					Longitude: db_trace.Longitude,
				},
			}

			chill.Traces = append(chill.Traces, trace)
		}

		db_photos, err := db.Photos(db.PhotoWhere.ChillID.EQ(db_chill.ID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db_photos.Select error:", err)

		}

		for _, db_photo := range db_photos {

			photo := &model.Photo{
				ID:        db_photo.ID,
				URL:       db_photo.URL,
				Timestamp: db_photo.Timestamp.Format("2006-01-02T15:04:05+09:00"),
			}

			chill.Photos = append(chill.Photos, photo)
		}
	}

	db_user_achievements, err := db.UserAchievements(db.UserAchievementWhere.UserID.EQ(uid)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_user_achievements.Select error:", err)
	}

	for _, db_user_achievement := range db_user_achievements {
		db_achievement, err := db.Achievements(db.AchievementWhere.ID.EQ(db_user_achievement.AchievementID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement.Select error:", err)

		}

		db_achievement_category, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(db_achievement.CategoryID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement_category.Select error:", err)

		}

		achievement := &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			Description: db_achievement.Description,
			DisplayName: db_achievement_category.DisplayName,
			Category: &model.AchievementCategory{
				ID:          db_achievement_category.ID,
				Name:        db_achievement_category.Name,
				DisplayName: db_achievement_category.DisplayName,
			},
		}

		result.Achievements = append(result.Achievements, achievement)
	}

	log.Println("db_user: ", db_user)

	err = db_user.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_user.Insert error:", err)
		return nil, err
	}

	return result, nil
}

func (u *UserService) GetUser(
	ctx context.Context,
	uid string,
) (*model.User, bool) {
	result := &model.User{}

	result.ID = uid

	db_user, err := db.Users(db.UserWhere.ID.EQ(uid)).One(ctx, u.Exec)
	if err != nil {
		log.Println("db_user.Select error:", err)
		return nil, false
	}

	result.Name = db_user.Name

	db_achievements, err := db.Achievements(db.AchievementWhere.ID.EQ(db_user.Avatar.String)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_achievement.Select error:", err)
		return nil, false
	}

	for _, db_achievement := range db_achievements {
		result.Avatar = &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			Description: db_achievement.Description,
			DisplayName: db_achievement.DisplayName,
		}
	
		db_achievement_category, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(db_achievement.CategoryID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement_category.Select error:", err)
			return nil, false
		}
	
		result.Avatar.Category = &model.AchievementCategory{
			ID:          db_achievement_category.ID,
			Name:        db_achievement_category.Name,
			DisplayName: db_achievement_category.DisplayName,
		}
	}

	db_user_chills, err := db.UserChills(db.UserChillWhere.UserID.EQ(uid)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_user_chills.Select error:", err)
		return nil, false
	}

	for _, db_user_chill := range db_user_chills {
		db_chill, err := db.Chills(db.ChillWhere.ID.EQ(db_user_chill.ChillID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_chill.Select error:", err)
			return nil, false
		}

		chill := &model.Chill{
			ID: db_chill.ID,
		}

		db_traces, err := db.TracePoints(db.TracePointWhere.ChillID.EQ(db_chill.ID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db_traces.Select error:", err)
			return nil, false
		}

		for _, db_trace := range db_traces {
			trace := &model.TracePoint{
				ID:        db_trace.ID,
				Timestamp: db_trace.Timestamp.Format("2006-01-02T15:04:05+09:00"),
				Coordinate: &model.Coordinate{
					Latitude:  db_trace.Latitude,
					Longitude: db_trace.Longitude,
				},
			}

			chill.Traces = append(chill.Traces, trace)
		}

		db_photos, err := db.Photos(db.PhotoWhere.ChillID.EQ(db_chill.ID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db_photos.Select error:", err)
			return nil, false
		}

		for _, db_photo := range db_photos {

			photo := &model.Photo{
				ID:        db_photo.ID,
				URL:       db_photo.URL,
				Timestamp: db_photo.Timestamp.Format("2006-01-02T15:04:05+09:00"),
			}

			chill.Photos = append(chill.Photos, photo)
		}
	}

	db_user_achievements, err := db.UserAchievements(db.UserAchievementWhere.UserID.EQ(uid)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_user_achievements.Select error:", err)
		return nil, false
	}

	for _, db_user_achievement := range db_user_achievements {
		db_achievement, err := db.Achievements(db.AchievementWhere.ID.EQ(db_user_achievement.AchievementID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement.Select error:", err)
			return nil, false
		}

		db_achievement_category, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(db_achievement.CategoryID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement_category.Select error:", err)
			return nil, false
		}

		achievement := &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			Description: db_achievement.Description,
			DisplayName: db_achievement_category.DisplayName,
			Category: &model.AchievementCategory{
				ID:          db_achievement_category.ID,
				Name:        db_achievement_category.Name,
				DisplayName: db_achievement_category.DisplayName,
			},
		}

		result.Achievements = append(result.Achievements, achievement)
	}

	return result, true
}

func (u *UserService) UpdateUser(
	ctx context.Context,
	user model.User,
	nameStr *string,
	avatarStr *string,
) (*model.User, error) {
	result := &model.User{}

	name := ""
	avatar := ""

	if nameStr != nil {
		name = *nameStr
	} else {
		name = user.Name
	}

	if avatarStr != nil {
		avatar = *avatarStr
	} else {
		avatar = user.Avatar.Name
	}

	db_user := &db.User{
		ID:     user.ID,
		Name:   name,
		Avatar: null.StringFromPtr(avatarStr),
	}

	result.ID = db_user.ID
	result.Name = db_user.Name

	db_achievements, err := db.Achievements(db.AchievementWhere.Name.EQ(avatar)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_achievement.Select error:", err)
		return nil, err
	}

	for _, db_achievement := range db_achievements {

		result.Avatar = &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			Description: db_achievement.Description,
			DisplayName: db_achievement.DisplayName,
		}

		db_user.Avatar = null.StringFrom(db_achievement.ID)
	
		db_achievement_category, err := db.AchievementCategories(db.AchievementCategoryWhere.ID.EQ(db_achievement.CategoryID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_achievement_category.Select error:", err)
			return nil, err
		}
	
		result.Avatar.Category = &model.AchievementCategory{
			ID:          db_achievement_category.ID,
			Name:        db_achievement_category.Name,
			DisplayName: db_achievement_category.DisplayName,
		}
	}

	db_user_chills, err := db.UserChills(db.UserChillWhere.UserID.EQ(user.ID)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_user_chills.Select error:", err)
		return nil, err
	}

	for _, db_user_chill := range db_user_chills {
		db_chill, err := db.Chills(db.ChillWhere.ID.EQ(db_user_chill.ChillID)).One(ctx, u.Exec)
		if err != nil {
			log.Println("db_chill.Select error:", err)
			return nil, err
		}

		chill := &model.Chill{
			ID: db_chill.ID,
		}

		db_traces, err := db.TracePoints(db.TracePointWhere.ChillID.EQ(db_chill.ID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db_traces.Select error:", err)
			return nil, err
		}

		for _, db_trace := range db_traces {
			trace := &model.TracePoint{
				ID:        db_trace.ID,
				Timestamp: db_trace.Timestamp.Format("2006-01-02T15:04:05+09:00"),
				Coordinate: &model.Coordinate{
					Latitude:  db_trace.Latitude,
					Longitude: db_trace.Longitude,
				},
			}

			chill.Traces = append(chill.Traces, trace)
		}

		db_photos, err := db.Photos(db.PhotoWhere.ChillID.EQ(db_chill.ID)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db_photos.Select error:", err)
			return nil, err
		}

		for _, db_photo := range db_photos {

			photo := &model.Photo{
				ID:        db_photo.ID,
				URL:       db_photo.URL,
				Timestamp: db_photo.Timestamp.Format("2006-01-02T15:04:05+09:00"),
			}

			chill.Photos = append(chill.Photos, photo)
		}
	}

	db_user_achievements, err := db.UserAchievements(db.UserAchievementWhere.UserID.EQ(user.ID)).All(ctx, u.Exec)
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

		achievement := &model.Achievement{
			ID:          db_achievement.ID,
			Name:        db_achievement.Name,
			Description: db_achievement.Description,
			DisplayName: db_achievement_category.DisplayName,
			Category: &model.AchievementCategory{
				ID:          db_achievement_category.ID,
				Name:        db_achievement_category.Name,
				DisplayName: db_achievement_category.DisplayName,
			},
		}

		result.Achievements = append(result.Achievements, achievement)
	}

	_, err = db_user.Update(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_user.Update error:", err)
		return nil, err
	}

	return result, nil

}
