package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userID string, expire time.Time, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    expire.Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, userID string, secret string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("Invalid token signature")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return fmt.Errorf("Invalid token body")
	}

	exp := time.Unix(int64(claims["exp"].(float64)), 0)

	if claims["userID"] == userID && exp.After(time.Now()) {
		return nil
	}

	return fmt.Errorf("Invalid or expired token")
}
