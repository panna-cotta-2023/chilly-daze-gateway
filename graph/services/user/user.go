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

	name := input.Name

	result := &model.User{}

	db_user := &db.User{
		ID:        uid,
		Name:      name,
		CreatedAt: time.Now(),
	}

	result.ID = uid
	result.Name = name

	err := db_user.Insert(ctx, u.Exec, boil.Infer())
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
	db_user, err := db.Users(db.UserWhere.ID.EQ(uid)).One(ctx, u.Exec)
	if err != nil {
		log.Println("db_user.Select error:", err)
		return nil, false
	}

	result := &model.User{
		ID:   db_user.ID,
		Name: db_user.Name,
		Avatar: &model.Achievement{
			ID: db_user.Avatar.String,
		},
	}

	return result, true
}

func (u *UserService) UpdateUser(
	ctx context.Context,
	userId string,
	input model.UpdateUserInput,
) (*model.User, error) {
	result := &model.User{}
	db_user, err := db.Users(db.UserWhere.ID.EQ(userId)).One(ctx, u.Exec)
	if err != nil {
		log.Println("db_user.Select error:", err)
		return nil, err
	}

	if input.Name != nil {
		result.Name = *input.Name
		db_user.Name = *input.Name
	}

	if input.Avatar != nil {
		dbUserAchievements, err := db.UserAchievements(db.UserAchievementWhere.UserID.EQ(userId)).All(ctx, u.Exec)
		if err != nil {
			log.Println("db.UserAchievements error:", err)
			return nil, err
		}

		for _, dbUserAchievement := range dbUserAchievements {
			dbAchievement, err := db.Achievements(db.AchievementWhere.ID.EQ(dbUserAchievement.AchievementID)).One(ctx, u.Exec)
			if err != nil {
				log.Println("db.Achievements error:", err)
				return nil, err
			}
			if dbAchievement.Name == db_user.Avatar.String {
				db_user.Avatar = null.StringFrom(dbAchievement.Name)
			}
		}
	}

	result.ID = db_user.ID

	_, err = db_user.Update(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_user.Update error:", err)
		return nil, err
	}

	return result, nil

}
