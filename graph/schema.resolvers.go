package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"goTicker/graph/generated"
	"goTicker/graph/model"
)

func (r *mutationResolver) PlaceOrder(ctx context.Context, input model.OrderCreate) (*model.OrderInfo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Orders(ctx context.Context) ([]*model.OrderInfo, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) placeOrder(ctx context.Context, input model.OrderCreate) (*model.OrderInfo, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *queryResolver) OrderInfo(ctx context.Context) ([]*model.OrderInfo, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *queryResolver) ServerStatus(ctx context.Context) ([]*model.ServerStatus, error) {
	panic(fmt.Errorf("not implemented"))
}
