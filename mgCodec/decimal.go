package mgCodec

import (
	"reflect"

	"github.com/d3v-friends/go-tools/fnError"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DecimalRegistry(registry *bsoncodec.Registry) *bsoncodec.Registry {
	var codec = &decimalValueCodec{}
	var typeof = reflect.TypeOf(decimal.Zero)

	registry.RegisterTypeEncoder(typeof, codec)
	registry.RegisterTypeDecoder(typeof, codec)

	return registry
}

const ErrInvalidDecimal = "invalid_decimal"

type decimalValueCodec struct{}

func (x *decimalValueCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, ok = value.Interface().(decimal.Decimal)
	if !ok {
		err = fnError.New(ErrInvalidDecimal)
		return
	}

	var dec primitive.Decimal128
	if dec, err = primitive.ParseDecimal128(i.String()); err != nil {
		return
	}

	if err = writer.WriteDecimal128(dec); err != nil {
		return
	}

	return
}

func (x *decimalValueCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	var d1 primitive.Decimal128
	if d1, err = reader.ReadDecimal128(); err != nil {
		return
	}

	var d2 decimal.Decimal
	if d2, err = decimal.NewFromString(d1.String()); err != nil {
		return
	}

	value.Set(reflect.ValueOf(d2))

	return
}
