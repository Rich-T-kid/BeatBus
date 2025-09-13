package internal

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var cfg = GetConfig()

var secretKey = []byte(cfg.JWTSecret)

type JWTHandler struct{}

func NewJWTHandler() *JWTHandler {
	return &JWTHandler{}
}
func (j *JWTHandler) CreateToken(username, roomID string, exp time.Duration) string {
	return issueHostToken(username, roomID, exp)
}
func (j *JWTHandler) VerifyToken(tokenString string) error {
	return verifyToken(tokenString)
}

func issueHostToken(username, roomID string, exp time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(exp).Unix(),
			"iss":      "BeatBus",
			"role":     "Host",
			"room_id":  roomID,
		})

	tokenString, _ := token.SignedString(secretKey)

	return tokenString
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

func RandomHash() string {
	return uuid.New().String()
}
