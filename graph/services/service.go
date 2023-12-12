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

type ChillService interface {
	AddChill(
		ctx context.Context,
		startChill model.StartChillInput,
		tracePoints []*model.TracePoint,
		photos []*model.Photo,
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
