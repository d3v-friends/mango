package typ

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	FnMigrate func(ctx context.Context, col *mongo.Collection) (memo string, err error)
)
