package m_migrate

import (
	"context"
	"github.com/d3v-friends/go-pure/fnReflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	DocMango struct {
		Id        primitive.ObjectID `bson:"_id"`
		ColNm     string             `bson:"colNm"`
		NextIdx   int                `bson:"nextIdx"`
		Histories []*DocMangoHistory `bson:"histories"`
		CreatedAt time.Time          `bson:"createdAt"`
	}

	DocMangoHistory struct {
		Memo      string    `bson:"memo"`
		CreatedAt time.Time `bson:"createdAt"`
	}
)

func (x *DocMango) GetMigrateList() FnMigrateList {
	return mgMango
}

func (x *DocMango) GetColNm() string {
	return docMangoNm
}

func (x *DocMango) GetID() primitive.ObjectID {
	return x.Id
}

/* ------------------------------------------------------------------------------------------------------------ */

const docMangoNm = "mango"

var mgMango = FnMigrateList{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "init indexing"
		_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{
						Key:   "colNm",
						Value: 1,
					},
				},
				Options: &options.IndexOptions{
					Unique: fnReflect.ToPointer(true),
				},
			},
		})

		return
	},
}
