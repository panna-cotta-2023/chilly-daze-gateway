package services

import (
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/chill"
	"chilly_daze_gateway/graph/services/photo"
	"chilly_daze_gateway/graph/services/trace"
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TraceService interface {
	AddTracePoints(
		ctx context.Context,
		input model.TracePointsInput,
	) ([]*model.TracePoint, error)
}

type PhotoService interface {
	AddPhotos(
		ctx context.Context,
		input model.PhotosInput,
	) ([]*model.Photo, error)
}

type ChillService interface {
	AddChill(
		ctx context.Context,
		startChill model.StartChillInput,
	) (*model.Chill, error)
}

type Services interface {
	TraceService
	PhotoService
	ChillService
}
type services struct {
	*trace.TraceService
	*photo.PhotoService
	*chill.ChillService
}

func New(exec boil.ContextExecutor) Services {
	return &services{
		TraceService: &trace.TraceService{Exec: exec},
		PhotoService: &photo.PhotoService{Exec: exec},
		ChillService: &chill.ChillService{Exec: exec},
	}
}
