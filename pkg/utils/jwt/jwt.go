package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type Payload struct {
	Key   string
	Value any
}

func WithPayload(key string, value any) Payload {
	return Payload{
		Key:   key,
		Value: value,
	}
}

func New(secret string, expire, iat int64, p ...Payload) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + expire
	claims["iat"] = iat

	for _, v := range p {
		claims[v.Key] = v.Value
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secret))
}

func Parse(token, secret string) (map[string]any, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetClaim(token, secret string, key string) (any, error) {
	claims, err := Parse(token, secret)
	if err != nil {
		return nil, err
	}

	value, exists := claims[key]
	if !exists {
		return nil, errors.New(fmt.Sprintf("claim %s not exists", key))
	}

	return value, nil
}
