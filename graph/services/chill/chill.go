package chill

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/lib"
	"context"
	"log"

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

	createTimeStamp, err := lib.ParseTimestamp(startChill.Timestamp)
	if err != nil {
		log.Println("lib.ParseTimestamp error:", err)
		return nil, err
	}

	result.Traces = append(result.Traces, &model.TracePoint{
		ID:        uuid.New().String(),
		Timestamp: createTimeStamp.Format("2006-01-02T15:04:05+09:00"),
		Coordinate: &model.Coordinate{
			Latitude:  startChill.Coordinate.Latitude,
			Longitude: startChill.Coordinate.Longitude,
		},
	})

	dbChill := &db.Chill{
		ID:        result.ID,
		CreatedAt: createTimeStamp,
	}

	err = dbChill.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("dbChill.Insert error:", err)
		return nil, err
	}

	dbTracePoint := &db.TracePoint{
		ID:        result.Traces[0].ID,
		Timestamp: createTimeStamp,
		ChillID:   result.ID,
		Latitude:  startChill.Coordinate.Latitude,
		Longitude: startChill.Coordinate.Longitude,
	}

	err = dbTracePoint.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("dbTracePoint.Insert error:", err)
		return nil, err
	}

	return result, nil
}

func (u *ChillService) EndChill(
	ctx context.Context,
	endChill model.EndChillInput,
	userId string,
) (*model.Chill, error) {
	result := &model.Chill{
		ID:     endChill.ID,
		Traces: []*model.TracePoint{},
		DistanceMeters: endChill.DistanceMeters,
	}

	for _, tracePoint := range endChill.TracePoints {
		timestamp, err := lib.ParseTimestamp(tracePoint.Timestamp)
		if err != nil {
			log.Println("lib.ParseTimestamp error:", err)
			return nil, err
		}

		dbTracePoint := &db.TracePoint{
			ID:        uuid.New().String(),
			Timestamp: timestamp,
			ChillID:   endChill.ID,
			Latitude:  tracePoint.Coordinate.Latitude,
			Longitude: tracePoint.Coordinate.Longitude,
		}

		err = dbTracePoint.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			log.Println("dbTracePoint.Insert error:", err)
			return nil, err
		}
	}
	
	if endChill.Photo != nil {
		photo := endChill.Photo

		timestamp, err := lib.ParseTimestamp(photo.Timestamp)
		if err != nil {
			log.Println("lib.ParseTimestamp error:", err)
			return nil, err
		}

		dbPhoto := &db.Photo{
			ID:        uuid.New().String(),
			ChillID:   endChill.ID,
			Timestamp: timestamp,
			URL:       photo.URL,
		}

		result.Photo = &model.Photo{
			ID:        dbPhoto.ChillID,
			Timestamp: timestamp.Format("2006-01-02T15:04:05+09:00"),
			URL:       dbPhoto.URL,
		}

		err = dbPhoto.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			log.Println("dbPhoto.Insert error:", err)
			return nil, err
		}
	}
	

	createTimeStamp, err := lib.ParseTimestamp(endChill.Timestamp)
	if err != nil {
		log.Println("lib.ParseTimestamp error:", err)
		return nil, err
	}

	endTimeStamp, err := lib.ParseTimestamp(endChill.Timestamp)
	if err != nil {
		log.Println("lib.ParseTimestamp error:", err)
		return nil, err
	}

	dbChill := &db.Chill{
		ID:        result.ID,
		CreatedAt: createTimeStamp,
		EndedAt:   null.TimeFrom(endTimeStamp),
	}

	_, err = dbChill.Update(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("dbChill.Update error:", err)
		return nil, err
	}

	return result, nil
}

func (u *ChillService) AddUserChill(
	ctx context.Context,
	userID string,
	chillID string,
) error {
	dbUserChill := &db.UserChill{
		UserID:  userID,
		ChillID: chillID,
	}

	err := dbUserChill.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("dbUserChill.Insert error:", err)
		return err
	}

	return nil
}

func (u *ChillService) GetChillsByUserId(
	ctx context.Context,
	userID string,
) ([]*model.Chill, error) {
	dbUserChills, err := db.UserChills(
		db.UserChillWhere.UserID.EQ(userID),
	).All(ctx, u.Exec)
	if err != nil {
		log.Println("db.UserChills error:", err)
		return nil, err
	}

	result := []*model.Chill{}

	for _, dbUserChill := range dbUserChills {
		dbChill, err := db.Chills(
			db.ChillWhere.ID.EQ(dbUserChill.ChillID),
		).One(ctx, u.Exec)
		if err != nil {
			log.Println("db.Chills error:", err)
			return nil, err
		}

		dbTracePoints, err := db.TracePoints(
			db.TracePointWhere.ChillID.EQ(dbChill.ID),
		).All(ctx, u.Exec)
		if err != nil {
			log.Println("db.TracePoints error:", err)
			return nil, err
		}

		traces := []*model.TracePoint{}

		for _, dbTracePoint := range dbTracePoints {
			traces = append(traces, &model.TracePoint{
				ID:        dbTracePoint.ID,
				Timestamp: dbTracePoint.Timestamp.Format("2006-01-02T15:04:05+09:00"),
				Coordinate: &model.Coordinate{
					Latitude:  dbTracePoint.Latitude,
					Longitude: dbTracePoint.Longitude,
				},
			})
		}

		dbPhotos, err := db.Photos(
			db.PhotoWhere.ChillID.EQ(dbChill.ID),
		).All(ctx, u.Exec)
		if err != nil {
			log.Println("db.Photos error:", err)
			return nil, err
		}

		photo := &model.Photo{}

		for _, dbPhoto := range dbPhotos {
			photo = &model.Photo{
				ID:        dbPhoto.ID,
				Timestamp: dbPhoto.Timestamp.Format("2006-01-02T15:04:05+09:00"),
				URL:       dbPhoto.URL,
			}
		}

		result = append(result, &model.Chill{
			ID:     dbChill.ID,
			Traces: traces,
			Photo: photo,
		})
	}

	return result, nil
}

