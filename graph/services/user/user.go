package user

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"context"
	"log"

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
	

	result := &model.User{
		ID:     uid,
		Name:   input.Name,
		Avatar: input.Avatar,
	}

	db_user := &db.User{
		ID:        result.ID,
		Name:      result.Name,
		AvatarURL: result.Avatar,
	}

	err := db_user.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_user.Insert error:", err)
	}

	return result, nil
}

func (u *UserService) GetUser(
	ctx context.Context,
	uid string,
) (*model.User, bool) {
	result := &model.User{}

	db_user, err := db.Users(db.UserWhere.ID.EQ(uid)).One(ctx, u.Exec)
	if err != nil {
		log.Println("db_user.Select error:", err)
		return nil, false
	}

	result.ID = db_user.ID
	result.Name = db_user.Name
	result.Avatar = db_user.AvatarURL
	
	return result, true
}

func (u *UserService) UpdateUserName(
	ctx context.Context,
	uid string,
	name string,
) (*model.User, error) {

	result := &model.User{}

	db_user, err := db.Users(db.UserWhere.ID.EQ(uid)).One(ctx, u.Exec)
	if err != nil {
		log.Println("db_user.Select error:", err)
		return nil, err
	}

	db_user.Name = name

	_, err = db_user.Update(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_user.Update error:", err)
		return nil, err
	}

	result.ID = db_user.ID
	result.Name = db_user.Name
	result.Avatar = db_user.AvatarURL

	return result, nil
}

func (u *UserService) UpdateUserAvatar(
	ctx context.Context,
	uid string,
	avatar string,
) (*model.User, error) {

	result := &model.User{}

	db_user, err := db.Users(db.UserWhere.ID.EQ(uid)).One(ctx, u.Exec)
	if err != nil {
		log.Println("db_user.Select error:", err)
		return nil, err
	}

	db_user.AvatarURL = avatar

	_, err = db_user.Update(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_user.Update error:", err)
		return nil, err
	}

	result.ID = db_user.ID
	result.Name = db_user.Name
	result.Avatar = db_user.AvatarURL

	return result, nil
}


