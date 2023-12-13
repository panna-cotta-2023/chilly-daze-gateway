package chill

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/lib"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
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

	createTimeStampString := lib.CovertTimestampString(startChill.Timestamp)

	createTimeStamp, err := time.Parse(time.RFC3339, createTimeStampString)
	if err != nil {
		log.Println("time.Parse error:", err)
		return nil, err
	}

	result.Traces = append(result.Traces, &model.TracePoint{
		ID:        uuid.New().String(),
		Timestamp: createTimeStampString,
		Coordinate: &model.Coordinate{
			Latitude:  startChill.Coordinate.Latitude,
			Longitude: startChill.Coordinate.Longitude,
		},
	})

	db_chill := &db.Chill{
		ID:        result.ID,
		CreatedAt: createTimeStamp,
	}

	err = db_chill.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_chill.Insert error:", err)
		return nil, err
	}

	db_tracePoint := &db.TracePoint{
		ID:        result.Traces[0].ID,
		Timestamp: createTimeStamp,
		ChillID:   result.ID,
		Latitude:  startChill.Coordinate.Latitude,
		Longitude: startChill.Coordinate.Longitude,
	}

	err = db_tracePoint.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_tracePoint.Insert error:", err)
		return nil, err
	}

	return result, nil
}

func (u *ChillService) EndChill(
	ctx context.Context,
	endChill model.EndChillInput,
) (*model.Chill, error) {
	result := &model.Chill{
		ID:     endChill.ID,
		Traces: []*model.TracePoint{},
	}

	result.Traces = append(result.Traces, &model.TracePoint{
		ID:        uuid.New().String(),
		Timestamp: endChill.Timestamp,
		Coordinate: &model.Coordinate{
			Latitude:  endChill.Coordinate.Latitude,
			Longitude: endChill.Coordinate.Longitude,
		},
	})

	createTimeStampString := lib.CovertTimestampString(endChill.Timestamp)

	createTimeStamp, err := time.Parse(time.RFC3339, createTimeStampString)
	if err != nil {
		log.Println("time.Parse error:", err)
		return nil, err
	}

	endTimeStampString := lib.CovertTimestampString(endChill.Timestamp)

	endTimeStamp, err := time.Parse(time.RFC3339, endTimeStampString)
	if err != nil {
		log.Println("time.Parse error:", err)
		return nil, err
	}

	db_chill := &db.Chill{
		ID:        result.ID,
		CreatedAt: createTimeStamp,
		EndedAt:   null.TimeFrom(endTimeStamp),
	}

	_, err = db_chill.Update(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_chill.Update error:", err)
		return nil, err
	}

	return result, nil
}

func (u *ChillService) AddUserChill(
	ctx context.Context,
	userID string,
	chillID string,
) error {
	db_userChill := &db.UserChill{
		UserID:  userID,
		ChillID: chillID,
	}

	err := db_userChill.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_userChill.Insert error:", err)
		return err
	}

	return nil
}

func (u *ChillService) GetChillsByUserId(
	ctx context.Context,
	userID string,
) ([]*model.Chill, error) {
	db_userChills, err := db.UserChills(
		db.UserChillWhere.UserID.EQ(userID),
	).All(ctx, u.Exec)
	if err != nil {
		log.Println("db.UserChills error:", err)
		return nil, err
	}

	result := []*model.Chill{}

	for _, db_userChill := range db_userChills {
		db_chill, err := db.Chills(
			db.ChillWhere.ID.EQ(db_userChill.ChillID),
		).One(ctx, u.Exec)
		if err != nil {
			log.Println("db.Chills error:", err)
			return nil, err
		}

		db_tracePoints, err := db.TracePoints(
			db.TracePointWhere.ChillID.EQ(db_chill.ID),
		).All(ctx, u.Exec)
		if err != nil {
			log.Println("db.TracePoints error:", err)
			return nil, err
		}

		traces := []*model.TracePoint{}

		for _, db_tracePoint := range db_tracePoints {
			traces = append(traces, &model.TracePoint{
				ID:        db_tracePoint.ID,
				Timestamp: db_tracePoint.Timestamp.Format("2006-01-02T15:04:05+09:00"),
				Coordinate: &model.Coordinate{
					Latitude:  db_tracePoint.Latitude,
					Longitude: db_tracePoint.Longitude,
				},
			})
		}

		result = append(result, &model.Chill{
			ID:        db_chill.ID,
			Traces:    traces,
			CreatedAt: db_chill.CreatedAt.Format("2006-01-02T15:04:05+09:00"),
			EndedAt:   db_chill.EndedAt.Time.Format("2006-01-02T15:04:05+09:00"),
		})
	}

	return result, nil

}
