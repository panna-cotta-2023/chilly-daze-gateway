package services

import (
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/chill"
	"chilly_daze_gateway/graph/services/photo"
	"chilly_daze_gateway/graph/services/trace"
	"chilly_daze_gateway/graph/services/user"
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserService interface {
	CreateUser(
		ctx context.Context,
		input model.RegisterUserInput,
	) (*model.User, error)
}

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
	StartChill(
		ctx context.Context,
		startChill model.StartChillInput,
	) (*model.Chill, error)
	EndChill(
		ctx context.Context,
		endChill model.EndChillInput,
	) (*model.Chill, error)
}

type Services interface {
	UserService
	TraceService
	PhotoService
	ChillService
}
type services struct {
	*user.UserService
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
