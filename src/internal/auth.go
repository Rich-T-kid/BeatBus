package internal

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var cfg = GetConfig()

var secretKey = []byte(cfg.JWTSecret)

type JWTHandler struct{}

func NewJWTHandler() *JWTHandler {
	return &JWTHandler{}
}
func (j *JWTHandler) CreateToken(username string, exp time.Duration) (string, error) {
	return " ", nil
}
func (j *JWTHandler) VerifyToken(tokenString string) error {
	return verifyToken(tokenString)
}

func issueHostToken(username, uuid, roomID string, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"uuid":     uuid,
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(exp).Unix(),
			"iss":      "BeatBus",
			"role":     "Host",
			"room_id":  roomID,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
