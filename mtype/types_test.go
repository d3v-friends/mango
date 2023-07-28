package mtype

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestTypes(test *testing.T) {

}

type (
	testModel struct {
		Id primitive.ObjectID
	}
)

func (x testModel) GetID() primitive.ObjectID {
	return x.Id
}
