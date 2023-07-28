package transact

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	testModel struct {
		Id        primitive.ObjectID `bson:"_id"`
		Name      string             `bson:"name"`
		InTrx     bool               `bson:"inTrx"`
		CreatedAt time.Time          `bson:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}
)

func (x *testModel) GetCollectionNm() string {
	return "transactModel"
}

func (x *testModel) GetMigrateList() minterface.FnMigrateList {
	return make(minterface.FnMigrateList, 0)
}

func (x *testModel) GetID() primitive.ObjectID {
	return x.Id
}

func (x *testModel) SetID(id primitive.ObjectID) {
	x.Id = id
}
