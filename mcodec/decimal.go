package mcodec

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type (
	DecimalCodec struct{}
	IfCodec      interface {
		bsoncodec.ValueEncoder
		bsoncodec.ValueDecoder
	}
)

func AppendDecimalCodec(registry *bsoncodec.Registry) *bsoncodec.Registry {
	var _, codec = DecimalCodecRegister()
	registry.RegisterTypeEncoder(reflect.TypeOf(decimal.Decimal{}), codec)
	registry.RegisterTypeDecoder(reflect.TypeOf(primitive.Decimal128{}), codec)
	return registry
}

var _ IfCodec = &DecimalCodec{}

func DecimalCodecRegister() (reflect.Type, IfCodec) {
	return reflect.TypeOf(decimal.Decimal{}), &DecimalCodec{}
}

func (dc *DecimalCodec) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	dec, ok := val.Interface().(decimal.Decimal)
	if !ok {
		return errors.New("invalidDecimal")
	}

	primDec, err := primitive.ParseDecimal128(dec.String())
	if err != nil {
		return errors.New("invalidDecimal")
	}
	return vw.WriteDecimal128(primDec)
}

func (dc *DecimalCodec) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	primDec, err := vr.ReadDecimal128()
	if err != nil {
		return errors.New("invalidDecimal")
	}

	dec, err := decimal.NewFromString(primDec.String())
	if err != nil {
		return errors.New("invalidDecimal")
	}

	val.Set(reflect.ValueOf(dec))
	return nil
}
