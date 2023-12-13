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

	result.Traces = append(result.Traces, &model.TracePoint{
		ID:        uuid.New().String(),
		Timestamp: startChill.Timestamp,
		Coordinate: &model.Coordinate{
			Latitude:  startChill.Coordinate.Latitude,
			Longitude: startChill.Coordinate.Longitude,
		},
	})

	createTimeStampString := lib.CovertTimestampString(startChill.Timestamp)

	createTimeStamp, err := time.Parse(time.RFC3339, createTimeStampString)
	if err != nil {
		log.Println("time.Parse error:", err)
		return nil, err
	}

	db_chill := &db.Chill{
		ID:        result.ID,
		CreatedAt: createTimeStamp,
	}

	err = db_chill.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_chill.Insert error:", err)
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
