package trace

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TraceService struct {
	Exec boil.ContextExecutor
}

func (u *TraceService) AddTracePoints(
	ctx context.Context,
	input model.TracePointsInput,
	chillID string,
) (*model.Chill, error) {
	chill := &model.Chill{}

	tracePoints := input.TracePoints

	for _, tracePoint := range tracePoints {

		timestamp, err := time.Parse(time.RFC3339, tracePoint.Timestamp)
		if err != nil {
			return nil, err
		}

		db_tracePoint := &db.TracePoint{
			ID:        uuid.New().String(),
			ChillID:   chillID,
			Latitude:  tracePoint.Coordinate.Latitude,
			Longitude: tracePoint.Coordinate.Longitude,
			Timestamp: timestamp,
		}

		err = db_tracePoint.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			return nil, err
		}
	}

	return chill, nil
}
