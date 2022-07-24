package token

import (
	"CP_Discussion/env"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var tokenErrs = []error{
	jwt.ErrTokenMalformed,
	jwt.ErrTokenExpired,
	jwt.ErrTokenNotValidYet,
}

type Claims struct {
	UserID string
	jwt.RegisteredClaims
}

func CreatToken(inUserID string) (string, error) {
	claims := Claims{
		UserID: inUserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(24) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(env.JWTKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(inToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(inToken, &Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(env.JWTKey), nil
		},
	)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			for _, tokenErr := range tokenErrs {
				if ve.Is(tokenErr) {
					return nil, tokenErr
				}
			}
		}
		return nil, errors.New("couldn't handle this token")
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}
