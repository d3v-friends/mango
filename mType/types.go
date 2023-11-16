package mType

type IndexModel struct {
	Key     map[string]int64 `bson:"key"`
	Name    string           `bson:"name"`
	Version int64            `bson:"v"`
}
