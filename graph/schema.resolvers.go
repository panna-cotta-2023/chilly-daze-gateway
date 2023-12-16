package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.41

import (
	"chilly_daze_gateway/graph/model"
	"context"
	"fmt"
	"log"
)

// Category is the resolver for the category field.
func (r *achievementResolver) Category(ctx context.Context, obj *model.Achievement) (*model.AchievementCategory, error) {
	return r.Srv.GetAchievementCategoryByAchievement(ctx, obj)
}

// Achievements is the resolver for the achievements field.
func (r *achievementCategoryResolver) Achievements(ctx context.Context, obj *model.AchievementCategory) ([]*model.Achievement, error) {
	return r.Srv.GetAchievementsByAchievementCategory(ctx, obj)
}

// Traces is the resolver for the traces field.
func (r *chillResolver) Traces(ctx context.Context, obj *model.Chill) ([]*model.TracePoint, error) {
	return r.Srv.GetTracePointsByChill(ctx, obj)
}

// Photo is the resolver for the photo field.
func (r *chillResolver) Photo(ctx context.Context, obj *model.Chill) (*model.Photo, error) {
	return r.Srv.GetPhotoByChill(ctx, obj)
}

// NewAchievements is the resolver for the newAchievements field.
func (r *chillResolver) NewAchievements(ctx context.Context, obj *model.Chill) ([]*model.Achievement, error) {
	userId := GetAuthToken(ctx)
	return r.Srv.GetNewAchievements(ctx, obj, userId)
}

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterUserInput) (*model.User, error) {
	userId := GetAuthToken(ctx)
	return r.Srv.CreateUser(ctx, input, userId)
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUserInput) (*model.User, error) {
	userId := GetAuthToken(ctx)
	return r.Srv.UpdateUser(ctx, userId, input)
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context) (*model.User, error) {
	userId := GetAuthToken(ctx)
	return r.Srv.DeleteUser(ctx, userId)
}

// StartChill is the resolver for the startChill field.
func (r *mutationResolver) StartChill(ctx context.Context, input model.StartChillInput) (*model.Chill, error) {
	userId := GetAuthToken(ctx)
	err := r.Srv.DeleteChillAfterOneDay(ctx)
	if err != nil {
		log.Println("r.Srv.DeleteChillAfterOneDay error:", err)
	}

	return r.Srv.StartChill(ctx, userId, input)
}

// EndChill is the resolver for the endChill field.
func (r *mutationResolver) EndChill(ctx context.Context, input model.EndChillInput) (*model.Chill, error) {
	userId := GetAuthToken(ctx)
	return r.Srv.EndChill(ctx, input, userId)
}

// AddTracePoints is the resolver for the AddTracePoints field.
func (r *mutationResolver) AddTracePoints(ctx context.Context, input model.TracePointsInput) ([]*model.TracePoint, error) {
	return r.Srv.AddTracePoints(ctx, input.TracePoints, input.ID)
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	userId := GetAuthToken(ctx)
	user, ok := r.Srv.GetUser(ctx, userId)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// Achievements is the resolver for the achievements field.
func (r *queryResolver) Achievements(ctx context.Context) ([]*model.Achievement, error) {
	return r.Srv.GetAchievements(ctx)
}

// AchievementCategories is the resolver for the achievementCategories field.
func (r *queryResolver) AchievementCategories(ctx context.Context) ([]*model.AchievementCategory, error) {
	return r.Srv.GetAchievementCategories(ctx)
}

// Avatar is the resolver for the avatar field.
func (r *userResolver) Avatar(ctx context.Context, obj *model.User) (*model.Achievement, error) {
	return r.Srv.GetAvatarByUser(ctx, obj)
}

// Chills is the resolver for the chills field.
func (r *userResolver) Chills(ctx context.Context, obj *model.User) ([]*model.Chill, error) {
	return r.Srv.GetChillsByUserId(ctx, obj.ID)
}

// Achievements is the resolver for the achievements field.
func (r *userResolver) Achievements(ctx context.Context, obj *model.User) ([]*model.Achievement, error) {
	return r.Srv.GetAchievementsByUserId(ctx, obj.ID)
}

// Achievement returns AchievementResolver implementation.
func (r *Resolver) Achievement() AchievementResolver { return &achievementResolver{r} }

// AchievementCategory returns AchievementCategoryResolver implementation.
func (r *Resolver) AchievementCategory() AchievementCategoryResolver {
	return &achievementCategoryResolver{r}
}

// Chill returns ChillResolver implementation.
func (r *Resolver) Chill() ChillResolver { return &chillResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type achievementResolver struct{ *Resolver }
type achievementCategoryResolver struct{ *Resolver }
type chillResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
