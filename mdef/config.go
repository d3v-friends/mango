package mdef

import (
	"fmt"
	"github.com/d3v-friends/go-pure/fnCases"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/gertd/go-pluralize"
	"regexp"
	"strconv"
	"strings"
)

var plural = pluralize.NewClient()

/* ------------------------------------------------------------------------------------------------------------ */

type Config struct {
	Version   Version                 `json:"version" yaml:"version"`
	Registry  map[MangoType]*Registry `json:"registry" yaml:"registry"`
	Documents map[DocNm]Document      `json:"documents" yaml:"documents"`
}

/* ------------------------------------------------------------------------------------------------------------ */

type Document struct {
	Field   map[DocFieldNm]DocField `json:"field" yaml:"field"`
	Times   []*TimeType             `json:"times" yaml:"times"`
	Indexes []*DocIndex             `json:"indexes" yaml:"indexes"`
}

/* ------------------------------------------------------------------------------------------------------------ */

type DocIndex struct {
	Key    DocIndexKey `json:"key" yaml:"key"`
	Unique bool        `json:"unique" yaml:"unique"`
}

/* ------------------------------------------------------------------------------------------------------------ */

type DocNm string

func (x DocNm) DocumentNm() string {
	return fnCases.CamelCase(plural.Plural(x.String()))
}

func (x DocNm) String() string {
	return string(x)
}

/* ------------------------------------------------------------------------------------------------------------ */

type DocFieldNm string

func (x DocFieldNm) String() string {
	return string(x)
}

func (x DocFieldNm) FieldNm() string {
	return fnCases.CamelCase(x.String())
}

/* ------------------------------------------------------------------------------------------------------------ */

type DocField any

// todo 이곳에 any 타입 reflect 이용하여 타입 정리하여 리턴하기

/* ------------------------------------------------------------------------------------------------------------ */

type DocFieldType string

const (
	DocFieldTypeMap    DocFieldType = "MAP"
	DocFieldTypeArray  DocFieldType = "ARRAY"
	DocFieldTypeObject DocFieldType = "OBJECT"
)

var DocFieldTypeAll = []DocFieldType{
	DocFieldTypeMap,
	DocFieldTypeArray,
	DocFieldTypeObject,
}

func (x DocFieldType) IsValid() bool {
	for _, fieldType := range DocFieldTypeAll {
		if fieldType == x {
			return true
		}
	}
	return false
}

/* ------------------------------------------------------------------------------------------------------------ */

type DocIndexKey [][]string

func (x DocIndexKey) IsValid() bool {
	for _, idxes := range x {
		if len(idxes) != 2 {
			return false
		}

		var order = Order(idxes[1])
		if !order.IsValid() {
			return false
		}
	}
	return true
}

/* ------------------------------------------------------------------------------------------------------------ */

type Order string

func (x Order) IsValid() bool {
	for _, order := range OrderAll {
		if x == order {
			return true
		}
	}
	return false
}

const (
	OrderASC  Order = "asc"
	OrderDESC Order = "desc"
)

var OrderAll = []Order{
	OrderASC,
	OrderDESC,
}

/* ------------------------------------------------------------------------------------------------------------ */

type Registry struct {
	Type string  `json:"type" yaml:"type"`
	Fn   *string `json:"fn" yaml:"fn"`
}

/* ------------------------------------------------------------------------------------------------------------ */

type Version string

var regexpVersion = fnPanic.OnValue(regexp.Compile(`^[0-9|.]+/g`))

func (x Version) IsValid(vers ...int) (err error) {
	if len(vers) == 0 {
		err = fmt.Errorf("invali version check: vers=%d", vers)
		return
	}

	var ls = strings.Split(regexpVersion.FindString(x.String()), ".")
	if len(ls) == 0 {
		err = fmt.Errorf("invalid version: version=%s", x.String())
		return
	}

	for i, ver := range vers {
		if len(ls) <= i {
			break
		}

		var checkVer int
		if checkVer, err = strconv.Atoi(ls[i]); err != nil {
			return
		}

		if ver != checkVer {
			err = fmt.Errorf("invalid version: version=%s, i=%d, ver=%d", x.String(), i, ver)
			return
		}
	}

	return
}

func (x Version) String() string {
	return string(x)
}

/* ------------------------------------------------------------------------------------------------------------ */

type MangoType string
type MangoTypeMap map[MangoType]*Registry

var mangoTypeMap = MangoTypeMap{
	"id": &Registry{
		Type: "go.mongodb.org/mongo-driver/bson/primitive.ObjectID",
	},
	"decimal": &Registry{
		Type: "go.mongodb.org/mongo-driver/bson/primitive.ObjectID",
	},
	"string": &Registry{
		Type: "string",
	},
	"int": &Registry{
		Type: "int64",
	},
	"time": &Registry{
		Type: "time.Time",
	},
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
