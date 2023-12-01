package helpers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

func CreateTokenJWT(id uint64) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    viper.GetString("JWT.ISSUER"),
		Subject:   strconv.Itoa(int(id)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	token, err := claims.SignedString([]byte(viper.GetString("JWT.SECRET")))
	if err != nil {
		return token, err
	}
	return token, nil
}

func DecodeTokenJWT(c *fiber.Ctx) (*jwt.Token, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid token")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("JWT.SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
