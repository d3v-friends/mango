package fnMango

import (
	"context"
	"github.com/d3v-friends/mango/migrate"
	"github.com/d3v-friends/mango/mtype"
	"go.mongodb.org/mongo-driver/mongo"
)

func Migrate(ctx context.Context, db *mongo.Database, models ...mtype.IfModel) (err error) {
	return migrate.V1(ctx, db, models...)
}
