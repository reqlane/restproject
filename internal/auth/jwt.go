package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userID int, username, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiresIn := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"uid":  userID,
		"user": username,
		"role": role,
	}

	if jwtExpiresIn != "" {
		duration, err := time.ParseDuration(jwtExpiresIn)
		if err != nil {
			return "", err
		}
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ParseToken(token string) (jwt.MapClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.New("token expired")
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New("token malformed")
		default:
			return nil, err
		}
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
