package token

import (
	"CP_Discussion/env"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ResetPWDClaims struct {
	Email string
	jwt.RegisteredClaims
}

func CreateResetPWDToken(inEmail string) (string, error) {
	claims := ResetPWDClaims{
		Email: inEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(10) * time.Minute)),
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

func ParseResetPWDToken(inToken string) (*ResetPWDClaims, error) {
	token, err := jwt.ParseWithClaims(inToken, &ResetPWDClaims{},
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
	if claims, ok := token.Claims.(*ResetPWDClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}
