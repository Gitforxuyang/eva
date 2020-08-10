package utils

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"reflect"
)

func StructToMap(v interface{}) map[string]interface{} {
	if IsNil(v) {
		return make(map[string]interface{})
	}
	return structs.Map(v)
}

func StructToJson(v interface{}) string {
	if v == nil {
		return ""
	}
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}

	return false
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func NotNil(v interface{}, name string) {
	if IsNil(v) {
		panic(fmt.Sprintf("%s不能为Nil", name))
	}
}
