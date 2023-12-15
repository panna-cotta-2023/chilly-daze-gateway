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

		db_tracePoint := &db.TracePoint{
			ID:        uuid.New().String(),
			Timestamp: timestamp,
			ChillID:   endChill.ID,
			Latitude:  tracePoint.Coordinate.Latitude,
			Longitude: tracePoint.Coordinate.Longitude,
		}

		err = db_tracePoint.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			log.Println("db_tracePoint.Insert error:", err)
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

		db_photo := &db.Photo{
			ID:        uuid.New().String(),
			ChillID:   endChill.ID,
			Timestamp: timestamp,
			URL:       photo.URL,
		}

		result.Photo = &model.Photo{
			ID:        db_photo.ChillID,
			Timestamp: timestamp.Format("2006-01-02T15:04:05+09:00"),
			URL:       db_photo.URL,
		}

		err = db_photo.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			log.Println("db_photo.Insert error:", err)
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

		db_photos, err := db.Photos(
			db.PhotoWhere.ChillID.EQ(db_chill.ID),
		).All(ctx, u.Exec)
		if err != nil {
			log.Println("db.Photos error:", err)
			return nil, err
		}

		photo := &model.Photo{}

		for _, db_photo := range db_photos {
			photo = &model.Photo{
				ID:        db_photo.ID,
				Timestamp: db_photo.Timestamp.Format("2006-01-02T15:04:05+09:00"),
				URL:       db_photo.URL,
			}
		}

		result = append(result, &model.Chill{
			ID:     db_chill.ID,
			Traces: traces,
			Photo: photo,
		})
	}

	return result, nil
}

