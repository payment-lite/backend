package authController

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"payment-gateway-lite/database"
	"payment-gateway-lite/database/models"
	"payment-gateway-lite/utils/helpers"
	"strings"
)

func GetUser(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		// decode token
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		//	}
		//	return []byte(viper.GetString("JWT.SECRET")), nil
		//})
		//if err != nil {
		//	return c.JSON(helpers.ErrorResponse{
		//		Success: false,
		//		Message: helpers.MessageUnauthorized,
		//		Errors:  err.Error(),
		//	})
		//}
		token, err := helpers.DecodeTokenJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ErrorResponse{
				Success: false,
				Message: helpers.MessageUnauthorized,
				Errors:  err.Error(),
			})
		}

		// get claims dari token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			sub := claims["sub"].(string)
			var user models.User

			err := database.DBConn.Preload("Team").First(&user, sub).Error
			if err != nil {
				return helpers.InternalServerError(c, err)
			}

			return c.JSON(&helpers.SuccessResponse{
				Success: true,
				Message: helpers.MessageSuccess,
				Data:    user,
			})
		}

	}

	// return unauthorized
	return c.JSON(helpers.ErrorResponse{
		Success: false,
		Message: helpers.MessageUnauthorized,
		Errors:  nil,
	})
}
