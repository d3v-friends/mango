package mgregistry

import (
	"reflect"

	"github.com/d3v-friends/go-tools/fnError"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func DecimalRegistry(registry *bson.Registry) *bson.Registry {

	var codec = &decimalValueCodec{}
	var typeof = reflect.TypeOf(decimal.Zero)

	registry.RegisterTypeEncoder(typeof, codec)
	registry.RegisterTypeDecoder(typeof, codec)

	return registry
}

const ErrInvalidDecimal = "invalid_decimal"

type decimalValueCodec struct{}

func (x *decimalValueCodec) EncodeValue(
	_ bson.EncodeContext,
	writer bson.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, ok = value.Interface().(decimal.Decimal)
	if !ok {
		err = fnError.New(ErrInvalidDecimal)
		return
	}

	var dec bson.Decimal128
	if dec, err = bson.ParseDecimal128(i.String()); err != nil {
		return
	}

	if err = writer.WriteDecimal128(dec); err != nil {
		return
	}

	return
}

func (x *decimalValueCodec) DecodeValue(
	_ bson.DecodeContext,
	reader bson.ValueReader,
	value reflect.Value,
) (err error) {
	var d1 bson.Decimal128
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
