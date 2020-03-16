package common

import (
	"reflect"
	"strings"
)

// 解析DB的相关数据
func GetModelInfo(table interface{}) (joinFields string , joinValues string , values []interface{}){
	t := reflect.TypeOf(table).Elem()
	v := reflect.ValueOf(table).Elem()
	var fields []string
	var tmpValues []string
	values = make([]interface{} , 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName != "" {
			fieldValue := v.Field(i)
			fields = append(fields, fieldName)
			var realValue interface{}
			switch fieldValue.Type().Kind().String(){
			case "string" :
				realValue = fieldValue.String()
			case "int" , "int64" :
				realValue = fieldValue.Int()
			case "float" , "float32" , "float64" :
				realValue = fieldValue.Float()
			}
			values = append(values , &realValue)
			tmpValues = append(tmpValues , "?")
		}
	}
	joinFields = strings.Join(fields , ",")
	joinValues = strings.Join(tmpValues , ",")
	return
}