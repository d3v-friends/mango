package mgCodec

import (
	"fmt"
	"github.com/d3v-friends/go-tools/fnError"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

const (
	ErrFailDecodeEnums = "fail_decode_enums"
	ErrFailEncodeEnums = "fail_encode_enums"
)

// GrpcEnumCodec
// GRPC 의 enum 을 string 문자열로 입력하게 해주는 codec
type GrpcEnumCodec[E GrpcEnum[E]] struct {
	strMap map[string]int32
	to     reflect.Type
}

func (x *GrpcEnumCodec[E]) strToT(str string) (res E, err error) {
	var i, has = x.strMap[str]
	if !has {
		err = fnError.NewFields(
			ErrFailDecodeEnums,
			map[string]any{
				"type":  reflect.TypeOf(new(E)).Name(),
				"value": str,
			},
		)
		return
	}

	res = (*new(E)).New(i)

	return
}

func (x *GrpcEnumCodec[E]) tToStr(t E) (res string, err error) {
	res = t.String()
	return
}

func (x *GrpcEnumCodec[E]) DecodeValue(
	_ bsoncodec.DecodeContext,
	vr bsonrw.ValueReader,
	val reflect.Value,
) (err error) {
	var str string
	if str, err = vr.ReadString(); err != nil {
		return
	}

	var res E
	if res, err = x.strToT(str); err != nil {
		return
	}

	val.Set(reflect.ValueOf(res))
	return
}

func (x *GrpcEnumCodec[E]) EncodeValue(
	_ bsoncodec.EncodeContext,
	vw bsonrw.ValueWriter,
	val reflect.Value,
) (err error) {
	var enum, isOk = val.Interface().(E)
	if !isOk {
		err = fnError.NewFields(
			ErrFailEncodeEnums,
			map[string]any{
				"type": reflect.TypeOf(new(E)).Name(),
			},
		)
		return
	}

	var str string
	if str, err = x.tToStr(enum); err != nil {
		return
	}

	return vw.WriteString(str)
}

func NewGrpcEnum[E GrpcEnum[E]](
	value map[string]int32,
) func(registry *bsoncodec.Registry) *bsoncodec.Registry {
	return func(registry *bsoncodec.Registry) *bsoncodec.Registry {
		var i = &GrpcEnumCodec[E]{
			strMap: value,
			to:     reflect.TypeOf(*new(E)),
		}
		registry.RegisterTypeEncoder(i.to, i)
		registry.RegisterTypeDecoder(i.to, i)
		return registry
	}
}

type GrpcEnum[T any] interface {
	fmt.Stringer
	New(int32) T
}
