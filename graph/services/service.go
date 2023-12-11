package services

import (
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/chill"
	"context"

	"database/sql"
)

type ChillService interface {
	StartChill(
		ctx context.Context,
		input model.StartChillInput,
	) (*model.Chill, error)
}

type Services interface {
	ChillService
}

type services struct {
	*chill.ChillService
}

func New(db *sql.DB) Services {
	return &services{
		ChillService: &chill.ChillService{
			Db: db,
		},
	}
}
