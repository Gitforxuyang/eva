package utils

import (
	"encoding/json"
	"github.com/fatih/structs"
)

func StructToMap(v interface{}) map[string]interface{} {

	return structs.Map(v)
}

func StructToJson(v interface{}) string {
	if v == nil {
		return ""
	}
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
