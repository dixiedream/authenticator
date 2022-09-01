package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var accessTokenSecret string = os.Getenv("JWT_ACCESS_SECRET")
var refreshTokenSecret string = os.Getenv("JWT_REFRESH_SECRET")

type Payload struct {
	Hostname string
	Role     int
}

type Claim struct {
	Name string
	Role int
	jwt.StandardClaims
}

func sign(payload *Payload, secret string, expiration int64) (string, error) {
	claims := Claim{Name: payload.Hostname, Role: payload.Role, StandardClaims: jwt.StandardClaims{ExpiresAt: expiration}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func isValid(t string, secret string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(t, Claim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing algorithm")
		}

		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(Claim); ok && token.Valid {
		return &Payload{Hostname: claims.Name, Role: claims.Role}, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, errors.New("Not a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("Token expired or not active")
		} else {
			log.Println(err)
			return nil, err
		}
	} else {
		return nil, err
	}
}

func GenerateAccessToken(payload *Payload) (string, error) {
	exp := time.Now().Add(time.Minute * 15).Unix()
	return sign(payload, accessTokenSecret, exp)
}

func AccessTokenIsValid(t string) (*Payload, error) {
	payload, err := isValid(t, accessTokenSecret)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func GenerateRefreshToken(payload *Payload) (string, error) {
	exp := time.Now().Add(time.Hour * 24 * 365).Unix()
	return sign(payload, refreshTokenSecret, exp)
}

func RefreshTokenIsValid(t string) (*Payload, error) {
	if t == "" {
		return nil, errors.New("Missing refresh token")
	}
	payload, err := isValid(t, refreshTokenSecret)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
