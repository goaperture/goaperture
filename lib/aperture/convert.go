package aperture

import (
	"net/url"
	"reflect"
	"strconv"
)

func getParamsToStruct[I any](values url.Values, to *I) {
	ref := reflect.ValueOf(to).Elem()

	for i := 0; i < ref.NumField(); i++ {
		key := ref.Type().Field(i).Name
		value := values.Get(ref.Type().Field(i).Tag.Get("json"))
		setField(to, key, value)
	}
}

func setField[T any](obj *T, fieldName string, value string) {
	v := reflect.ValueOf(obj).Elem()
	field := v.FieldByName(fieldName)

	if field.IsValid() && field.CanSet() {
		switch field.Type().Name() {
		case "int":
			intval, _ := strconv.ParseInt(value, 10, 64)
			field.SetInt(intval)
		case "float":
			floatval, _ := strconv.ParseFloat(value, 64)
			field.SetFloat(floatval)
		case "bool":
			field.SetBool(value == "1")
		case "string":
			field.SetString(value)
		}
	}
}
