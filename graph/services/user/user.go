package user

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserService struct {
	Exec boil.ContextExecutor
}

func (u *UserService) CreateUser(
	ctx context.Context,
	input model.RegisterUserInput,
) (*model.User, error) {
	result := &model.User{
		ID:     uuid.New().String(),
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
