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

func (u *PhotoService) AddPhoto(
	ctx context.Context,
	input *model.PhotoInput,
	chillId string,
) (*model.Photo, error) {

	if input == nil {
		return &model.Photo{}, nil
	}

	result := &model.Photo{}

	photo := input

	timestamp, err := lib.ParseTimestamp(photo.Timestamp)
	if err != nil {
		log.Println("lib.ParseTimestamp error:", err)
		return nil, err
	}

	db_photo := &db.Photo{
		ID:        uuid.New().String(),
		ChillID:   chillId,
		Timestamp: timestamp,
		URL:       photo.URL,
	}

	result = &model.Photo{
		ID:        db_photo.ChillID,
		Timestamp: timestamp.Format("2006-01-02T15:04:05+09:00"),
		URL:       db_photo.URL,
	}

	err = db_photo.Insert(ctx, u.Exec, boil.Infer())
	if err != nil {
		log.Println("db_photo.Insert error:", err)
		return nil, err
	}

	return result, nil
}

func (u *PhotoService) GetPhotosByChill(
	ctx context.Context,
	chill *model.Chill,
) ([]*model.Photo, error) {
	result := []*model.Photo{}

	db_photos, err := db.Photos(db.PhotoWhere.ChillID.EQ(chill.ID)).All(ctx, u.Exec)
	if err != nil {
		log.Println("db_photos.Select error:", err)
		return nil, err
	}

	for _, db_photo := range db_photos {

		photo := &model.Photo{
			ID:        db_photo.ID,
			URL:       db_photo.URL,
			Timestamp: db_photo.Timestamp.Format("2006-01-02T15:04:05+09:00"),
		}

		result = append(result, photo)
	}

	return result, nil
}
