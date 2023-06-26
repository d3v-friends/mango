package cores

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ManagedModel[T any] struct {
		Data  T       `bson:"inline"`
		Mango *FMango `bson:"__mango"`
	}

	FMango struct {
		IsLock bool `bson:"isLock"`
	}

	IfManagedModel interface {
		GetID() primitive.ObjectID
		CollectionNm() string
	}

	DocMongoIndex struct {
		Key        map[string]int `bson:"key"`
		Name       string         `bson:"name"`
		V          int            `bson:"v"`
		Background *bool          `bson:"background"`
		Unique     *bool          `bson:"unique"`
	}

	DocMongoIndexList []DocMongoIndex
)

const (
	idxManagedModelIsLock = "__mango.isLock_1"
)

// RegisterManagedModel 관리할 모델에 기본 데이터와 인덱싱을 추가해준다
func RegisterManagedModel(
	ctx context.Context,
	db *mongo.Database,
	modelList ...IfManagedModel,
) (err error) {
	for _, model := range modelList {
		col := db.Collection(model.CollectionNm())
		var cur *mongo.Cursor
		if cur, err = col.Indexes().List(ctx); err != nil {
			return
		}

		ls := make(DocMongoIndexList, 0)
		if err = cur.All(ctx, &ls); err != nil {
			return
		}

		if ls.HasByName(idxManagedModelIsLock) {

		} else {

		}
	}

	return
}

func (x *DocMongoIndexList) HasByName(name string) bool {
	for i := range *x {
		if (*x)[i].Name == name {
			return true
		}
	}
	return false
}
