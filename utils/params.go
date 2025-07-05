package utils

import (
	"net/url"
	"reflect"
	"strconv"
)

func BuildQueryParams(filter interface{}) string {
	values := url.Values{}
	v := reflect.ValueOf(filter)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ""
	}
	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("json")

		if tag == "" || tag == "-" {
			continue
		}

		value := ""
		switch field.Kind() {
		case reflect.String:
			value = field.String()
		case reflect.Int, reflect.Int64:
			if field.Int() != 0 {
				value = strconv.FormatInt(field.Int(), 10)
			}
		}

		if value != "" {
			values.Add(tag, value)
		}
	}
	return values.Encode()
}
