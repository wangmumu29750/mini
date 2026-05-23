package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Principal struct {
	UserID   uint64 `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type claims struct {
	UserID   uint64 `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(secret string, ttl time.Duration, principal Principal) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID:   principal.UserID,
		Username: principal.Username,
		Role:     principal.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", principal.UserID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
		},
	})
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, raw string) (Principal, error) {
	token, err := jwt.ParseWithClaims(raw, &claims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return Principal{}, err
	}

	claimsValue, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return Principal{}, fmt.Errorf("invalid token")
	}

	return Principal{
		UserID:   claimsValue.UserID,
		Username: claimsValue.Username,
		Role:     claimsValue.Role,
	}, nil
}
