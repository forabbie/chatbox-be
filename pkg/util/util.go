package util

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Names(structure interface{}) []string {
	list := []string{}

	names := reflect.ValueOf(structure)

	if names.Kind() == reflect.Ptr {
		names = names.Elem()
	}

	if names.Kind() == reflect.Struct {
		for i := 0; i < names.NumField(); i++ {
			list = append(list, names.Type().Field(i).Name)
		}
	}

	return list
}

func Search(haystack []string, needle string, sensitive bool) int {
	for key, val := range haystack {
		switch sensitive {
		case true:
			if val == needle {
				return key
			}
		case false:
			if strings.ToLower(val) == needle {
				return key
			}
		}
	}

	return -1
}

func GenerateCode(length int) string {
	code := make([]string, length)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len(code); i++ {
		code[i] = strconv.Itoa(rand.Intn(9))
	}

	return strings.Join(code, "")
}

func ReplacePlaceholders(query string, argCount int) (string, error) {
	var sb strings.Builder
	argIndex := 1
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			if argIndex > argCount {
				return "", fmt.Errorf("more placeholders than arguments")
			}
			sb.WriteString(fmt.Sprintf("$%d", argIndex))
			argIndex++
		} else {
			sb.WriteByte(query[i])
		}
	}
	return sb.String(), nil
}
