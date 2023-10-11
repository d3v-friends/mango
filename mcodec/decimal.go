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
	registry.RegisterTypeEncoder(DecimalCodecRegister())
	registry.RegisterTypeDecoder(DecimalCodecRegister())
	return registry
}

var _ IfCodec = &DecimalCodec{}

func DecimalCodecRegister() (reflect.Type, IfCodec) {
	return reflect.TypeOf(decimal.Decimal{}), &DecimalCodec{}
}

func (dc *DecimalCodec) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) (err error) {
	var dec, ok = val.Interface().(decimal.Decimal)
	if !ok {
		err = errors.New("invalid decimal")
		return
	}

	var primitiveDecimal primitive.Decimal128
	if primitiveDecimal, err = primitive.ParseDecimal128(dec.String()); err != nil {
		return
	}

	return vw.WriteDecimal128(primitiveDecimal)
}

func (dc *DecimalCodec) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) (err error) {
	var primitiveDecimal primitive.Decimal128
	if primitiveDecimal, err = vr.ReadDecimal128(); err != nil {
		return errors.New("invalid decimal")
	}

	var dec decimal.Decimal
	if dec, err = decimal.NewFromString(primitiveDecimal.String()); err != nil {
		return errors.New("invalid decimal")
	}

	val.Set(reflect.ValueOf(dec))

	return
}
