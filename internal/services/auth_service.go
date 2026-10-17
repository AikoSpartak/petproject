package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

const ttlAccessToken = 24 * time.Hour
const ttlRefreshToken = 3 * 24 * time.Hour

func getJWT() ([]byte, error) {
	key, ok := os.LookupEnv("JWT_KEY")

	if !ok {
		return nil, fmt.Errorf("invalid secret key")
	}

	return []byte(key), nil
}

func getRefreshToken() ([]byte, error) {

	key, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		return nil, fmt.Errorf("invalid secret key")
	}
	return []byte(key), nil
}

func GenerateAccessToken(userID, role string) (string, error) {

	secretKey, err := getJWT()
	if err != nil {
		return "", fmt.Errorf("failed to get JWT")
	}

	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttlAccessToken)),
			Subject:   "access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			key, err := getJWT()
			if err != nil {
				return nil, fmt.Errorf("failed to get JWT key: %w", err)
			}
			return key, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GenerateRefreshToken(userID, role string) (string, error) {
	secretKey, err := getRefreshToken()
	if err != nil {
		return "", fmt.Errorf("failed to get refresh key: %w", err)
	}

	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttlRefreshToken)),
			Subject:   "refresh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
