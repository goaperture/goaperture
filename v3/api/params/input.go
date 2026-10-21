package params

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func GetInput[I any](r *http.Request) I {
	var input I

	switch r.Method {
	case http.MethodGet:
		toStruct(r.URL.Query(), &input)
	default:
		json.NewDecoder(r.Body).Decode(&input)
	}

	return input
}

func toStruct[I any](values url.Values, to *I) {
	ref := reflect.ValueOf(to).Elem()
	t := ref.Type()

	for i := 0; i < ref.NumField(); i++ {
		field := t.Field(i)
		fieldVal := ref.Field(i)

		// Берем имя из тега json (игнорируя параметры вроде ,string)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		cleanTag := strings.Split(tag, ",")[0]

		val := values.Get(cleanTag)
		if val == "" {
			continue
		}

		// Логика конвертации
		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(val)
		case reflect.Int, reflect.Int64:
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				fieldVal.SetInt(i)
			}
		case reflect.Bool:
			if b, err := strconv.ParseBool(val); err == nil {
				fieldVal.SetBool(b)
			}
			// Добавьте другие типы по необходимости
		}
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
