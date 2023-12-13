package trace

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/lib"
	"context"
	"log"
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
) ([]*model.TracePoint, error) {
	tracePoints := input.TracePoints

	result := []*model.TracePoint{}

	for _, tracePoint := range tracePoints {
		timestampString := lib.CovertTimestampString(tracePoint.Timestamp)

		timestamp, err := time.Parse(time.RFC3339, timestampString)
		if err != nil {
			log.Println("time.Parse error:", err)
			return nil, err
		}

		db_tracePoint := &db.TracePoint{
			ID:        uuid.New().String(),
			ChillID:   input.ID,
			Latitude:  tracePoint.Coordinate.Latitude,
			Longitude: tracePoint.Coordinate.Longitude,
			Timestamp: timestamp,
		}

		result = append(result, &model.TracePoint{
			ID:        db_tracePoint.ChillID,
			Timestamp: db_tracePoint.Timestamp.Format("2006-01-02T15:04:05+09:00"),
			Coordinate: &model.Coordinate{
				Latitude:  db_tracePoint.Latitude,
				Longitude: db_tracePoint.Longitude,
			},
		})

		err = db_tracePoint.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			log.Println("db_tracePoint.Insert error:", err)
			return nil, err
		}
	}

	return result, nil
}

func (u *TraceService) GetTracesByChillId(
	ctx context.Context,
	chillId string,
) ([]*model.TracePoint, error) {

	db_traces, err := db.TracePoints(
		db.TracePointWhere.ChillID.EQ(chillId),
	).All(ctx, u.Exec)
	if err != nil {
		log.Println("db.TracePoints error:", err)
		return nil, err
	}

	result := []*model.TracePoint{}

	for _, db_trace := range db_traces {
		result = append(result, &model.TracePoint{
			ID:        db_trace.ID,
			Timestamp: db_trace.Timestamp.Format("2006-01-02T15:04:05+09:00"),
			Coordinate: &model.Coordinate{
				Latitude:  db_trace.Latitude,
				Longitude: db_trace.Longitude,
			},
		})
	}

	return result, nil
}
