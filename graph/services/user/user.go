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

	dbUser := &db.User{
		ID:        uid,
		Name:      name,
		CreatedAt: time.Now(),
	}

	result.ID = uid
	result.Name = name

	err := dbUser.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("dbUser.Insert error:", err)
		return nil, err
	}

	return result, nil
}

func (u *UserService) GetUser(
	ctx context.Context,
	uid string,
) (*model.User, bool) {
	dbUser, err := db.Users(db.UserWhere.ID.EQ(uid)).One(ctx, u.Exec)
	if err != nil {
		log.Println("dbUser.Select error:", err)
		return nil, false
	}

	result := &model.User{
		ID:   dbUser.ID,
		Name: dbUser.Name,
		Avatar: &model.Achievement{
			ID: dbUser.Avatar.String,
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
	dbUser, err := db.Users(db.UserWhere.ID.EQ(userId)).One(ctx, u.Exec)
	if err != nil {
		log.Println("dbUser.Select error:", err)
		return nil, err
	}

	if input.Name != nil {
		result.Name = *input.Name
		dbUser.Name = *input.Name
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
			if dbAchievement.Name == dbUser.Avatar.String {
				dbUser.Avatar = null.StringFrom(dbAchievement.Name)
			}
		}
	}

	result.ID = dbUser.ID

	_, err = dbUser.Update(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("dbUser.Update error:", err)
		return nil, err
	}

	return result, nil

}
