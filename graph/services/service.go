package services

import (
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/chill"

	"database/sql"
)

type ChillService interface {
	StartChill() (*model.Chill, error)
}

type Services interface{
	ChillService
}

type services struct{
	*chill.ChillService
}

func New(db *sql.DB) Services {
	return &services{
		ChillService: &chill.ChillService{
			Db: db,
		},
	}
}
