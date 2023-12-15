package services

import (
	"chilly_daze_gateway/graph/model"
	"chilly_daze_gateway/graph/services/achievement"
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
		uid string,
	) (*model.User, error)
	GetUser(
		ctx context.Context,
		uid string,
	) (*model.User, bool)
	UpdateUser(
		ctx context.Context,
		user model.User,
		nameStr *string,
		avatarStr *string,
	) (*model.User, error)
}

type TraceService interface {
	AddTracePoints(
		ctx context.Context,
		input model.TracePointsInput,
	) ([]*model.TracePoint, error)
	GetTracePointsByChill(
		ctx context.Context,
		chill *model.Chill,
	) ([]*model.TracePoint, error)
}

type PhotoService interface {
	AddPhotos(
		ctx context.Context,
		input model.PhotosInput,
	) ([]*model.Photo, error)
	GetPhotosByChill(
		ctx context.Context,
		chill *model.Chill,
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
		userId string,
	) (*model.Chill, error)
	AddUserChill(
		ctx context.Context,
		userID string,
		chillID string,
	) error
	GetChillsByUserId(
		ctx context.Context,
		userID string,
	) ([]*model.Chill, error)
}

type AchievementService interface {
	GetAchievementsByUserId(
		ctx context.Context,
		user_id string,
	) ([]*model.Achievement, error)
	GetAchievements(
		ctx context.Context,
	) ([]*model.Achievement, error)
	GetAchievementCategories(
		ctx context.Context,
	) ([]*model.AchievementCategory, error)
	AddChillAchievement(
		ctx context.Context,
		user_id string,
		chill_id string,
		achievement_id []string,
	) error
	GetAvatarByUser(
		ctx context.Context,
		user *model.User,
	) (*model.Achievement, error)
	GetAchievementCategoryByAchievement(
		ctx context.Context,
		achievement *model.Achievement,
	) (*model.AchievementCategory, error)
	GetAchievementsByAchievementCategory(
		ctx context.Context,
		achievementCategory *model.AchievementCategory,
	) ([]*model.Achievement, error)
}

type Services interface {
	UserService
	TraceService
	PhotoService
	ChillService
	AchievementService
}
type services struct {
	*user.UserService
	*trace.TraceService
	*photo.PhotoService
	*chill.ChillService
	*achievement.AchievementService
}

func New(exec boil.ContextExecutor) Services {
	return &services{
		UserService:  &user.UserService{Exec: exec},
		TraceService: &trace.TraceService{Exec: exec},
		PhotoService: &photo.PhotoService{Exec: exec},
		ChillService: &chill.ChillService{Exec: exec},
		AchievementService: &achievement.AchievementService{Exec: exec},
	}
}
