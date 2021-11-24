package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	mySigningKey = "WOW,MuchShibe,ToDogge"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

const (
	tokenExpiresIn = 15000
)

func main() {
	createdToken, err := ExampleNew([]byte(mySigningKey))
	if err != nil {
		fmt.Println("Creating token failed")
	}
	ExampleParse(createdToken, mySigningKey)
}

func ExampleNew(mySigningKey []byte) (string, error) {
	createdAt := time.Now().Unix()
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   "subject",
			ExpiresAt: createdAt + tokenExpiresIn,
			IssuedAt:  createdAt,
			NotBefore: createdAt,
		},
		Email: "email@gmail.com",
	})
	// Set some claims

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(mySigningKey)

	return tokenString, err
}

func ExampleParse(myToken string, myKey string) (string, string, error) {
	fmt.Println(myToken)
	parsed, err := jwt.ParseWithClaims("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImVtYWlsQGdtYWlsLmNvbSIsImV4cCI6MTYyMTk1NzU5NCwiaWF0IjoxNjIxOTQyNTk0LCJuYmYiOjE2MjE5NDI1OTQsInN1YiI6InN1YmplY3QifQ.iu7esx6-GEId3F0TOmcquQz-_jCWTMIAoOagvyb4ytI", &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(myKey), nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		fmt.Println(claims.Email)
		fmt.Println(claims.StandardClaims.Subject)
		return claims.StandardClaims.Subject, claims.Email, nil
	}
	return "", "", errors.New("invalid token")
}
