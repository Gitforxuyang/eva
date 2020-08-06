package utils

import (
	"encoding/json"
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
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
