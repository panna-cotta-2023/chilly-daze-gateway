package photo

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/lib"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type PhotoService struct {
	Exec boil.ContextExecutor
}

func (u *PhotoService) AddPhotos(
	ctx context.Context,
	input model.PhotosInput,
) ([]*model.Photo, error) {
	photos := input.Photos

	result := []*model.Photo{}

	for _, photo := range photos {

		timestamp, err := lib.ParseTimestamp(photo.Timestamp)

		db_photo := &db.Photo{
			ID:        uuid.New().String(),
			ChillID:   input.ID,
			Timestamp: timestamp,
			URL:       photo.URL,
		}

		result = append(result, &model.Photo{
			ID:        db_photo.ChillID,
			Timestamp: timestamp.Format("2006-01-02T15:04:05+09:00"),
			URL:       db_photo.URL,
		})

		err = db_photo.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			log.Println("db_photo.Insert error:", err)
			return nil, err
		}
	}

	return result, nil
}

func (u *PhotoService) GetPhotosByChillId(
	ctx context.Context,
	chillId string,
) ([]*model.Photo, error) {
	db_photos, err := db.Photos(
		db.PhotoWhere.ChillID.EQ(chillId),
	).All(ctx, u.Exec)
	if err != nil {
		log.Println("db.Photos error:", err)
		return nil, err
	}

	result := []*model.Photo{}

	for _, db_photo := range db_photos {
		result = append(result, &model.Photo{
			ID:        db_photo.ChillID,
			Timestamp: db_photo.Timestamp.Format("2006-01-02T15:04:05+09:00"),
			URL:       db_photo.URL,
		})
	}

	return result, nil
}
