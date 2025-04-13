package util

import (
	"math/rand"
	"reflect"
	"strings"
	"strconv"
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
