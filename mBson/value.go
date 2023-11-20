package mBson

import (
	"fmt"
	"github.com/d3v-friends/go-pure/fnLogger"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strconv"
)

//
//const bsonTag = "bson"
//
//type Value interface {
//}
//
//func UpdateBson(v any) (res bson.D, err error) {
//	var typeOf = reflect.TypeOf(v)
//	var valueOf = reflect.ValueOf(v)
//
//	var kind = typeOf.Kind()
//	switch kind {
//	case reflect.Pointer:
//		if valueOf.IsNil() {
//			res = make(bson.D, 0)
//			return
//		}
//		return UpdateBson(v)
//	case reflect.Struct:
//		return updateBson(valueOf, reflect.TypeOf(v), "")
//	default:
//		err = fmt.Errorf("invalid kind: kind=%s", kind.String())
//		return
//	}
//
//	panic("not impl")
//}
//
//func getValue(i any, tag string) (ls bson.D, err error) {
//	ls = make(bson.D, 0)
//	var v = reflect.ValueOf(i)
//
//	switch v.Kind() {
//	case reflect.Pointer:
//		if !v.CanInterface() {
//			break
//		}
//
//		var items bson.D
//		if items, err = getValue(v.Elem().Interface(), tag); err != nil {
//			return
//		}
//		ls = append(ls, items...)
//	case reflect.Struct:
//		var t = reflect.TypeOf(i)
//		for i := 0; i < v.NumField(); i++ {
//			var field = v.Field(i)
//			if !field.CanInterface() {
//				continue
//			}
//
//			var fieldTag = t.Field(i).Tag.Get("bson")
//			var items bson.D
//			if items, err = getValue(field.Interface(), mergeTag(tag, fieldTag)); err != nil {
//				return
//			}
//			ls = append(ls, items...)
//		}
//	case reflect.Invalid,
//		reflect.Chan,
//		reflect.Func,
//		reflect.Interface:
//		return
//	case reflect.Slice, reflect.Array:
//		ls = append(ls, bson.E{
//			Key:   tag,
//			Value: v,
//		})
//	case reflect.Map:
//		for _, key := range v.MapKeys() {
//			if !v.MapIndex(key).CanInterface() {
//				continue
//			}
//
//			var t = mergeTag(tag, key.String())
//			var items bson.D
//			if items, err = getValue(v.)
//
//		}
//		return
//	default:
//		ls = append(ls, bson.E{
//			Key:   tag,
//			Value: i,
//		})
//		return
//	}
//}
//func mergeTag(a, b string) (res string) {
//	if a == "" {
//		return b
//	}
//
//	if b == "" {
//		fnLogger.Panic("b is empty")
//	}
//
//	res = fmt.Sprintf("%s.%s", a, b)
//	return
//}
