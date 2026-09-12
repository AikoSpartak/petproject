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

func getJWT() ([]byte, error) {
	key, ok := os.LookupEnv("JWT_Key")
	if !ok {
		return nil, fmt.Errorf("invalid secret key")
	}

	return []byte(key), nil
}

func getRefreshToken() ([]byte, error) {

	key, ok := os.LookupEnv("REFRESH_KEY")
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
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



func GenerateRefreshToken(userID string) (string, error) {

	refreshToken, err := getRefreshToken()
	if err != nil {
		return "", fmt.Errorf("failed to get JWT")
	}

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * 24 * time.Hour)),
		Subject:   "refresh",
		ID:        userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(refreshToken)
}
