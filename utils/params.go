package utils

import (
	"net/url"
	"reflect"
	"strconv"

	"github.com/YHVCorp/vercelgo/schemas"
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

func BuildProjectDomainsParams(teamId string, opts *schemas.Options) string {
	values := url.Values{}
	values.Add("teamId", teamId)

	if opts != nil {
		if opts.Production != nil {
			values.Add("production", *opts.Production)
		}
		if opts.Target != nil {
			values.Add("target", *opts.Target)
		}
		if opts.CustomEnvironmentID != nil {
			values.Add("customEnvironmentId", *opts.CustomEnvironmentID)
		}
		if opts.GitBranch != nil {
			values.Add("gitBranch", *opts.GitBranch)
		}
		if opts.Redirects != nil {
			values.Add("redirects", *opts.Redirects)
		}
		if opts.Redirect != nil {
			values.Add("redirect", *opts.Redirect)
		}
		if opts.Verified != nil {
			values.Add("verified", *opts.Verified)
		}
		if opts.Limit != nil {
			values.Add("limit", strconv.Itoa(*opts.Limit))
		}
		if opts.Since != nil {
			values.Add("since", strconv.FormatInt(*opts.Since, 10))
		}
		if opts.Until != nil {
			values.Add("until", strconv.FormatInt(*opts.Until, 10))
		}
		if opts.Order != nil {
			values.Add("order", *opts.Order)
		}
		if opts.Slug != nil {
			values.Add("slug", *opts.Slug)
		}
	}

	return values.Encode()
}
