package authentication

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

const secretKey = "badforprod"

func GenerateToken(email string, userid int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":   email,
		"user_id": userid,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (claims jwt.MapClaims, err error) {
	log.Printf("Received token: %v\n", token)
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("could not parse token")
	}
	if !parsed.Valid {
		return claims, errors.New("invalid token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return claims, errors.New("invalid token claims")
	}
	exp := int64(claims["exp"].(float64))
	if exp < time.Now().Unix() {
		return claims, errors.New("token has expired")
	}
	return claims, err
}
