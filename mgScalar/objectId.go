package mgScalar

import (
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/d3v-friends/go-tools/fnError"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ErrInvalidObjectID = "invalid_object_id"

func MarshalObjectID(v primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Quote(v.Hex())))
	})
}

func UnmarshalObjectID(v any) (res primitive.ObjectID, err error) {
	switch t := v.(type) {
	case string:
		return primitive.ObjectIDFromHex(t)
	case []byte:
		return primitive.ObjectIDFromHex(string(t))
	default:
		err = fnError.NewF(ErrInvalidObjectID)
		return
	}
}
