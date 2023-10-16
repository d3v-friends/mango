package mango

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IfDoc interface {
	GetID() primitive.ObjectID
}

type Doc struct {
	ID        primitive.ObjectID `bson:"_id"`
	IsLock    bool               `bson:"isLock"`
	CreatedAt *time.Time         `bson:"createdAt"`
	UpdatedAt *time.Time         `bson:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt"`

	// privates
	isLoaded bool `bson:"-"`
}

/* ------------------------------------------------------------------------------------------------------------ */

type TimeType string

func (x TimeType) IsValid() bool {
	for _, timeType := range TimeTypeAll {
		if x == timeType {
			return true
		}
	}
	return false
}

const (
	TimeTypeCreatedAt TimeType = "CREATED_AT"
	TimeTypeUpdatedAt TimeType = "UPDATED_AT"
	TimeTypeDeletedAt TimeType = "DELETED_AT"
)

var TimeTypeAll = []TimeType{
	TimeTypeCreatedAt,
	TimeTypeUpdatedAt,
	TimeTypeDeletedAt,
}

/* ------------------------------------------------------------------------------------------------------------ */

type IfDocData interface {
	HookTime() []TimeType
}
