package services

import (
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/photo"
	"chilly_daze_gateway/graph/services/trace"
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TraceService interface {
	AddTracePoints(
		ctx context.Context,
		input model.TracePointsInput,
		chillID string,
	) ([]*model.TracePoint, error)
}

type PhotoService interface {
	AddPhotos(
		ctx context.Context,
		input model.PhotosInput,
		chillID string,
	) ([]*model.Photo, error)
}

type Services interface {
	TraceService
	PhotoService
}
type services struct {
	*trace.TraceService
	*photo.PhotoService
}

func New(exec boil.ContextExecutor) Services {
	return &services{
		TraceService: &trace.TraceService{Exec: exec},
		PhotoService: &photo.PhotoService{Exec: exec},
	}
}
