package utils

import (
	"fmt"
	"reflect"
	"sort"
)

type SignSort string

const (
	SortDesc SignSort = "desc"
	SortAsc  SignSort = "asc"
)

func GetSignValue(obj interface{}, sortDesc SignSort) string {
	kvs := make([]Kv, 0)
	object := reflect.TypeOf(obj)
	switch object.Kind() {
	case reflect.Struct:
		v := reflect.ValueOf(obj)
		for i := 0; i < object.NumField(); i++ {
			field := object.Field(i)
			if field.Tag.Get("json") != "-" && field.Tag.Get("json") != "" {
				kvs = append(kvs, Kv{Key: field.Tag.Get("json"), Val: getValue(v.Field(i), sortDesc)})
			}
		}
		if sortDesc == "desc" {
			sort.Sort(KvSliceDesc(kvs))
		} else {
			sort.Sort(KvSliceAsc(kvs))
		}
		var str string
		for _, v := range kvs {
			str = str + v.Val
		}
		return str
	case reflect.Pointer:
		myRef := object.Elem()
		objectVal := reflect.ValueOf(obj).Elem()
		for i := 0; i < myRef.NumField(); i++ {
			field := myRef.Field(i)
			if field.Tag.Get("json") != "-" && field.Tag.Get("json") != "" {
				kvs = append(kvs, Kv{Key: field.Tag.Get("json"), Val: getValue(objectVal.Field(i), sortDesc)})
			}
		}
		if sortDesc == "desc" {
			sort.Sort(KvSliceDesc(kvs))
		} else {
			sort.Sort(KvSliceAsc(kvs))
		}
		var str string
		for _, v := range kvs {
			str = str + v.Val
		}
		return str
	case reflect.Array:
		str := ""
		v := reflect.ValueOf(obj)
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			str += GetSignValue(element, sortDesc)
		}
		return str
	case reflect.Slice:
		str := ""
		v := reflect.ValueOf(obj)
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			str += GetSignValue(element.Interface(), sortDesc)
		}
		return str
	case reflect.Map:
		p := obj.(map[string]interface{})

		for k, v := range p {
			t := reflect.ValueOf(v)
			kvs = append(kvs, Kv{
				Key: k,
				Val: getValue(t, sortDesc),
			})
		}
		if sortDesc == "desc" {
			sort.Sort(KvSliceDesc(kvs))
		} else {
			sort.Sort(KvSliceAsc(kvs))
		}
		var str string
		for _, v := range kvs {
			str = str + v.Val
		}
		return str
	default:
		return ""
	}
}

func getValue(val reflect.Value, sort SignSort) string {
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprint(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprint(val.Uint())
	case reflect.Float64, reflect.Float32:
		return fmt.Sprint(val.Float())
	case reflect.Struct:
		return GetSignValue(val.Interface(), sort)
	case reflect.Array:
		str := ""
		for i := 0; i < val.Len(); i++ {
			element := val.Index(i)
			str += GetSignValue(element, sort)
		}
		return str
	case reflect.Slice:
		str := ""
		for i := 0; i < val.Len(); i++ {
			element := val.Index(i)
			str += GetSignValue(element.Interface(), sort)
		}
		return str
	default:
		return ""
	}
}
