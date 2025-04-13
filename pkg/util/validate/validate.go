package validate

import (
	"reflect"
	"regexp"
	"strings"
)

const (
	EmailAddressRegExpPattern = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	EmptyRegExpPattern = `^\s*$`
)

type Map map[string]interface{}

var (
	EmailAddressRegExp = regexp.MustCompile(EmailAddressRegExpPattern)
	EmptyRegExp = regexp.MustCompile(EmptyRegExpPattern)
)

func One(name string, field interface{}, validate string) Map {
	invalid := Map{}

	one := reflect.ValueOf(field)

	if reflect.TypeOf(field) == reflect.TypeOf(reflect.Value{}) {
		one = field.(reflect.Value)
	}

	if one.Kind() == reflect.Ptr {
		one = one.Elem()
	}

	Loop:
		for _, val := range strings.Split(validate, ",") {
			switch val {
			case "required":
				if !one.IsValid() {
					invalid = Map{ "name": name, "validate": val }

					break Loop
				}

				if EmptyRegExp.MatchString(one.String()) {
					invalid = Map{ "name": name, "validate": val }

					break Loop
				}

				if one.Kind() == reflect.Struct {
					if one.IsZero() {
						invalid = Map{ "name": name, "validate": val }

						break Loop
					}
				}
			case "emailaddress":
				if !EmailAddressRegExp.MatchString(one.String()) {
					invalid = Map{ "name": name, "validate": val }

					break Loop
				}
			}
		}

	return invalid
}

func All(structure interface{}) []Map {
	invalids := []Map{}

	all := reflect.ValueOf(structure)

	if all.Kind() == reflect.Ptr {
		all = all.Elem()
	}

	if all.Kind() == reflect.Struct {
		for i := 0; i < all.NumField(); i++ {
			name := strings.ToLower(all.Type().Field(i).Name)

			field := all.Field(i)

			validate := all.Type().Field(i).Tag.Get("validate")

			invalid := One(name, field, validate)

			if len(invalid) > 0 {
				invalids = append(invalids, invalid)
			}
		}
	}

	return invalids
}
