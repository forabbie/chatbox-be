package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strings"
)

func HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func HexEncode(b []byte) string {
	return hex.EncodeToString(b)
}

func Decrypt(value, key string) (string, error) {
	decodedKey, _ := HexDecode(key)

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", err
	}

	decodedValue, _ := HexDecode(value)

	if len(decodedValue) < aes.BlockSize {
		return "", errors.New("encrypted value is not valid")
	}

	ciphertext := decodedValue[aes.BlockSize:]

	plaintext := make([]byte, len(ciphertext))

	iv := decodedValue[:aes.BlockSize]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}

func Encrypt(value, key string) (string, error) {
	decodedKey, _ := HexDecode(key)

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", err
	}

	plaintext := []byte(value)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return HexEncode(ciphertext), nil
}

func RandomCode(length int) (string, error) {
	code := make([]string, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			return "", err
		}

		code[i] = n.String()
	}

	return strings.Join(code, ""), nil
}

func RandomHexCode(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func RandomBase64Code(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func Names(structure interface{}) []string {
	list := []string{}

	names := reflect.ValueOf(structure)

	if names.Kind() == reflect.Ptr {
		names = names.Elem()
	}

	if names.Kind() == reflect.Struct {
		for i := 0; i < names.NumField(); i++ {
			field := names.Type().Field(i)

			if tag := field.Tag.Get("json"); tag != "" {
				list = append(list, strings.Split(tag, ",")[0])
			} else {
				list = append(list, names.Type().Field(i).Name)
			}
		}
	}

	return list
}

func Find(haystack []int, needle int) int {
	for key, val := range haystack {
		if val == needle {
			return key
		}
	}

	return -1
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

func ConvertToJSON(data interface{}) interface{} {
	switch value := data.(type) {
	case map[interface{}]interface{}:
		out := map[string]interface{}{}

		for key, val := range value {
			out[fmt.Sprintf("%v", key)] = ConvertToJSON(val)
		}

		return out
	case []interface{}:
		for i, val := range value {
			value[i] = ConvertToJSON(val)
		}

		return value
	default:
		return value
	}
}

func IsNil(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsNil()
}

func Extract(text, field string) string {
	prefix := field + "="

	for _, val := range strings.Split(text, ", ") {
		if strings.HasPrefix(val, prefix) {
			return val[len(prefix):]
		}

		if val == field {
			return val
		}
	}

	return ""
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
