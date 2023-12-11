package chill

import (
	"chilly_daze_gateway/graph/model"
	"context"
	"database/sql"
)

type ChillService struct {
	Db *sql.DB
}

func (u *ChillService) StartChill(
	ctx context.Context,
	input model.StartChillInput,
) (*model.Chill, error) {
	
}
