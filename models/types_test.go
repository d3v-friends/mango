package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestTypes(test *testing.T) {
	// Value 리시버에서 작동한다
	filter := NewTrxOne[testModel]("testModel", &bson.M{
		"_id": primitive.NewObjectID(),
	})

	filter.GetCollectionName()

}

type (
	testModel struct {
		Id primitive.ObjectID
	}
)

func (x testModel) GetID() primitive.ObjectID {
	return x.Id
}
