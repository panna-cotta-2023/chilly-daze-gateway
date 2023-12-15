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
	user model.User,
	input model.UpdateUserInput,
) (*model.User, error) {
	result := &model.User{}

	name := ""

	if input.Name != nil {
		name = *input.Name
	} else {
		name = user.Name
	}

	db_user := &db.User{
		ID:     user.ID,
		Name:   name,
		Avatar: null.StringFromPtr(input.Avatar),
	}

	result.ID = db_user.ID
	result.Name = db_user.Name

	_, err := db_user.Update(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_user.Update error:", err)
		return nil, err
	}

	return result, nil

}
