package services

import (
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/trace"
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TraceService interface {
	AddTracePoints(
		ctx context.Context,
		input model.TracePointsInput,
		chillID string,
	) (*model.Chill, error)
}

type Services interface {
	TraceService
}
type services struct {
	*trace.TraceService
}

func New(exec boil.ContextExecutor) Services {
	return &services{
		TraceService: &trace.TraceService{Exec: exec},
	}
}
