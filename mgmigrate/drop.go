package mgmigrate

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Drop(
	db *mongo.Database,
	models ...MigratedModel,
) (err error) {
	var ctx = context.TODO()

	for _, model := range models {
		if err = db.
			Collection(model.GetColNm()).
			Drop(ctx); err != nil {
			return
		}

		if _, err = db.
			Collection(ColNm).
			DeleteOne(
				ctx,
				bson.M{
					FieldColNm: model.GetColNm(),
				},
			); err != nil {
			return
		}
	}

	return
}
