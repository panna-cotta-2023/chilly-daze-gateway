package photo

import (
	"chilly_daze_gateway/graph/db"
	"chilly_daze_gateway/graph/model"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type PhotoService struct {
	Exec boil.ContextExecutor
}

func (u *PhotoService) AddPhotos(
	ctx context.Context,
	input model.PhotosInput,
	chillID string,
) ([]*model.Photo, error) {
	photos := input.Photos

	result := []*model.Photo{}

	for _, photo := range photos {

		timestamp, err := time.Parse(time.RFC3339, photo.Timestamp)
		if err != nil {
			return nil, err
		}

		db_photo := &db.Photo{
			ID:        uuid.New().String(),
			ChillID:   chillID,
			Timestamp: timestamp,
			URL:       photo.URL,
		}

		result = append(result, &model.Photo{
			ID:        db_photo.ID,
			Timestamp: db_photo.Timestamp.Format("2006-01-02T15:04:05.00:00+00:00"),
			URL:       db_photo.URL,
		})

		err = db_photo.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
