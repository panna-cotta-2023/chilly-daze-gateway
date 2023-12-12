package chill

import (
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/db"
	"context"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ChillService struct {
	Exec boil.ContextExecutor
}

func (u *ChillService) AddChill(
	ctx context.Context,
	startChill model.StartChillInput,
) (*model.Chill, error) {
	result := &model.Chill{
		ID: uuid.New().String(),
	}

	db_chill := &db.Chill{
		ID: result.ID,
	}

	err := db_chill.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		return nil, err
	}

	return result, nil
}
