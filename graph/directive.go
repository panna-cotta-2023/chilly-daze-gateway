package graph

import (
	"chilly_daze_gateway/middleware/auth"
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/99designs/gqlgen/graphql"
)

var Directive DirectiveRoot = DirectiveRoot{
	IsAuthenticated: IsAuthenticated,
}

func IsAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalln("firebase.NewApp error:", err)
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalln("app.Auth error:", err)
		return nil, err
	}

	if idToken, ok := auth.GetIdToken(ctx); ok {
		_, err = client.VerifyIDToken(ctx, idToken)
		if err != nil {
			log.Println("client.VerifyIDToken error:", err)
			return nil, err
		}
	}

	return next(ctx)
}
