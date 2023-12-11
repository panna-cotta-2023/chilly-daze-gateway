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
) (*model.Chill, error) {
	chill := &model.Chill{}

	photos := input.Photos

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

		err = db_photo.Insert(ctx, u.Exec, boil.Infer())
		if err != nil {
			return nil, err
		}
	}

	return chill, nil
}
