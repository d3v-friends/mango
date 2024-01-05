package mType

import (
	"fmt"

	"github.com/d3v-friends/go-pure/fnParams"
)

type IndexModel struct {
	Key     map[string]int64 `bson:"key"`
	Name    string           `bson:"name"`
	Version int64            `bson:"v"`
}

type IndexModels []*IndexModel

func (x IndexModels) Has(key string, prefixes ...string) (has bool) {
	var prefix = fnParams.Get(prefixes)
	if prefix != "" {
		prefix = fmt.Sprintf("%s.", prefix)
	}

	var idKey = fmt.Sprintf("%s%s", prefix, key)
	for _, id := range x {
		if id.Name == idKey {
			return true
		}
	}
	return false
}
