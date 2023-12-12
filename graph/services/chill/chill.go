package chill

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ChillService struct {
	Exec boil.ContextExecutor
}

func (u *ChillService) StartChill(
	ctx context.Context,
	startChill model.StartChillInput,
) (*model.Chill, error) {
	result := &model.Chill{
		ID:     uuid.New().String(),
		Traces: []*model.TracePoint{},
	}

	result.Traces = append(result.Traces, &model.TracePoint{
		ID:        uuid.New().String(),
		Timestamp: startChill.Timestamp,
		Coordinate: &model.Coordinate{
			Latitude:  startChill.Coordinate.Latitude,
			Longitude: startChill.Coordinate.Longitude,
		},
	})

	createTimeStamp, err := time.Parse(time.RFC3339, startChill.Timestamp)
	if err != nil {
		return nil, err
	}

	db_chill := &db.Chill{
		ID:         result.ID,
		CreatedAt: createTimeStamp,
	}

	err = db_chill.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		return nil, err
	}

	return result, nil
}
