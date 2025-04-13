package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"
	"strings"

	jwtv4 "github.com/golang-jwt/jwt/v4"
)

const (
	KeyLength int = 32

	AuthScheme string = "Bearer"
)

var AccessTokenKey, RefreshTokenKey string

func GenerateKey(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func Auth(auth string, scheme string) string {
	auths := strings.Split(auth, scheme)

	if len(auths) == 2 {
		return strings.TrimSpace(auths[1])
	}

	return auth
}

func NewToken(sub interface{}, exp time.Duration, jti interface{}, key string) (string, error) {
	claims := jwtv4.MapClaims{
		// "iss": "",
		"sub": sub,
		// "aud": []string{},
		"exp": time.Now().Add(exp).Unix(),
		// "nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		// "jti": "",
	}

	if jti != nil {
		claims["jti"] = jti
	}

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func ParseToken(auth string, key string) (jwtv4.MapClaims, error) {
	token, err := jwtv4.ParseWithClaims(auth, jwtv4.MapClaims{}, func(token *jwtv4.Token) (interface{}, error) { return []byte(key), nil })
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwtv4.MapClaims)

	return claims, nil
}
