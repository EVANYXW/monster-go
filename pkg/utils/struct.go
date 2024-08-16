/**
 * @api post utils.
 *
 * User: yunshengzhu
 * Date: 2022/3/29
 * Time: 下午2:57
 */
package utils

import (
	"fmt"
	"reflect"
	"sort"
)

type Kv struct {
	Key string
	Val string
}

type KvSliceAsc []Kv

func (s KvSliceAsc) Len() int      { return len(s) }
func (s KvSliceAsc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s KvSliceAsc) Less(i, j int) bool { return s[i].Key < s[j].Key }

type KvSliceDesc []Kv

func (s KvSliceDesc) Len() int      { return len(s) }
func (s KvSliceDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s KvSliceDesc) Less(i, j int) bool { return s[i].Key > s[j].Key }

func GetSortAscStr(obj interface{}) string {
	kvs := make([]Kv, 0)
	object := reflect.TypeOf(obj)
	objectVal := reflect.ValueOf(obj).Elem()
	myRef := object.Elem()
	for i := 0; i < myRef.NumField(); i++ {
		field := myRef.Field(i)
		if field.Tag.Get("json") != "-" {
			kvs = append(kvs, Kv{Key: field.Tag.Get("json"), Val: GetValue(objectVal.Field(i))})
		}
	}
	sort.Sort(KvSliceAsc(kvs))
	fmt.Println(fmt.Sprintf("============:%+v", kvs))
	var str string
	for _, v := range kvs {
		str = str + v.Val
	}
	return str
}

func GetSortDescStr(obj interface{}) string {
	kvs := make([]Kv, 0)
	object := reflect.TypeOf(obj)
	objectVal := reflect.ValueOf(obj).Elem()
	myRef := object.Elem()
	for i := 0; i < myRef.NumField(); i++ {
		field := myRef.Field(i)
		if field.Tag.Get("json") != "-" {
			kvs = append(kvs, Kv{Key: field.Tag.Get("json"), Val: GetValue(objectVal.Field(i))})
		}
	}
	sort.Sort(KvSliceDesc(kvs))
	var str string
	for _, v := range kvs {
		str = str + v.Val
	}
	return str
}

func GetValue(val reflect.Value) string {
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprint(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprint(val.Uint())
	case reflect.Float64, reflect.Float32:
		return fmt.Sprint(val.Float())
	default:
		return ""
	}
}
