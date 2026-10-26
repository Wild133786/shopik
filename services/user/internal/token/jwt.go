package token

import (
	"shopik/services/user/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Role string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

type JWTIssuer struct {
	secret []byte
	issuer string
	TTL    time.Duration
}

func NewJWTIssuer(secret []byte, TTL time.Duration) (*JWTIssuer, error) {
	if len(secret) < 32 {
		return nil, domain.ErrInvalidJWTSecret
	}
	return &JWTIssuer{
		secret: secret,
		issuer: "shopik-user-service",
		TTL:    TTL,
	}, nil
}

func (i *JWTIssuer) Issuer(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    i.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(i.TTL * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(i.secret)
}
