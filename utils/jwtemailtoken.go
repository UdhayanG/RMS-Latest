package utils

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	tokenExpiresIn = 15000
)

// TokenHandler is the struct that generates JWT tokens.
type TokenHandler struct {
	PrivateKey *rsa.PrivateKey
}

// Claims is the struct that represent JWT claims according to https://tools.ietf.org/html/rfc7519#section-4
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// NewToken creates a token for email verification
// with the given subject and email.
func (e *TokenHandler) NewToken(subject string, email string) (string, error) {

	createdAt := time.Now().Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			ExpiresAt: createdAt + tokenExpiresIn,
			IssuedAt:  createdAt,
			NotBefore: createdAt,
		},
		Email: email,
	})
	return t.SignedString(e.PrivateKey)
}

// Verify validates a given token and returns subject and email form it.
func (e TokenHandler) Verify(token string) (string, string, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return &e.PrivateKey.PublicKey, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return claims.StandardClaims.Subject, claims.Email, nil
	}
	return "", "", errors.New("invalid token")
}
