package authController

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"payment-gateway-lite/database"
	"payment-gateway-lite/database/models"
	"payment-gateway-lite/utils/helpers"
)

type SignupRequest struct {
	Name     string `json:"name" validate:"required,gte=3,lte=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=50"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=50"`
}

type GoogleOauthRequest struct {
	Name  string `json:"name" validate:"required,gte=3,lte=100"`
	Email string `json:"email" validate:"required,email"`
}

var validate = validator.New()

func Signup(c *fiber.Ctx) error {
	var reqData SignupRequest
	//parse data dari client
	if err := c.BodyParser(&reqData); err != nil {
		return helpers.InternalServerError(c, err)
	}

	// Validate data dari client
	err := validate.Struct(reqData)
	if err != nil {
		return helpers.ErrorValidation(c, err)
	}

	// hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.InternalServerError(c, err)
	}

	// create user & return error jika user sudah ada
	var user models.User
	err = database.DBConn.Where("email = ?", reqData.Email).First(&user).Error
	if err == nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&helpers.ErrorResponse{
			Success: false,
			Message: "Email already exists",
			Errors:  nil,
		})
	}

	// create user
	user.Email = reqData.Email
	user.Name = reqData.Name
	user.Password = string(hashedPassword)
	err = database.DBConn.Create(&user).Error
	if err != nil {
		return helpers.InternalServerError(c, err)
	}
	// create team
	team := models.Team{
		OwnerID: &user.ID,
		Name:    user.Name,
	}
	err = database.DBConn.Create(&team).Error

	// set user team
	if err := database.DBConn.Model(&user).Update("team_id", team.OwnerID).Error; err != nil {
		return helpers.InternalServerError(c, err)
	}

	// create token
	token, err := helpers.CreateTokenJWT(user.ID)
	if err != nil {
		return helpers.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(&helpers.SuccessResponse{
		Success: true,
		Message: helpers.MessageSuccess,
		Data: fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

func Signin(c *fiber.Ctx) error {
	var reqData SigninRequest
	if err := c.BodyParser(&reqData); err != nil {
		return helpers.InternalServerError(c, err)
	}

	// Validate
	err := validate.Struct(reqData)
	if err != nil {
		return helpers.ErrorValidation(c, err)
	}

	// Find user by email & return error jika user tidak ditemukan
	var user models.User
	result := database.DBConn.Model(&models.User{}).First(&user, "email = ?", reqData.Email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "failed",
				"message": helpers.MessageUnauthorized,
			})
		}
	}

	// create token
	token, err := helpers.CreateTokenJWT(user.ID)
	if err != nil {
		return helpers.InternalServerError(c, err)
	}
	return c.JSON(helpers.SuccessResponse{
		Success: true,
		Message: "Signin success",
		Data: fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

// oauth
func GoogleOauthLogin(c *fiber.Ctx) error {
	var reqData GoogleOauthRequest

	err := c.BodyParser(&reqData)
	if err != nil {
		return helpers.InternalServerError(c, err)
	}

	// Validate
	err = validate.Struct(reqData)
	if err != nil {
		return helpers.ErrorValidation(c, err)
	}

	//random password
	randPassword, err := helpers.GenerateRandomPassword(20)
	if err != nil {
		return helpers.InternalServerError(c, err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randPassword), bcrypt.DefaultCost)
	if err != nil {
		return helpers.InternalServerError(c, err)
	}

	// create user jika belum ada
	var user models.User
	// cari user
	err = database.DBConn.Model(&models.User{}).First(&user, "email = ?", reqData.Email).Error
	// jika ga di temukan buat baru
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var team models.Team
			user = models.User{
				Name:     reqData.Name,
				Email:    reqData.Email,
				Password: string(hashedPassword),
			}
			if err := database.DBConn.Model(&models.User{}).Create(&user).Error; err != nil {
				return helpers.InternalServerError(c, err)
			}

			// create team
			team = models.Team{
				Name:    reqData.Name,
				OwnerID: &user.ID,
			}
			if err := database.DBConn.Model(&models.Team{}).Create(&team).Error; err != nil {
				return helpers.InternalServerError(c, err)
			}

			// Set User's TeamID
			if err := database.DBConn.Model(user).Update("team_id", team.OwnerID).Error; err != nil {
				return helpers.InternalServerError(c, err)
			}

			token, err := helpers.CreateTokenJWT(user.ID)
			if err != nil {
				return helpers.InternalServerError(c, err)
			}
			return c.JSON(helpers.SuccessResponse{
				Success: true,
				Message: "Signin success",
				Data: fiber.Map{
					"user":  user,
					"token": token,
				},
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(&helpers.ErrorResponse{
			Success: false,
			Message: helpers.MessageUnauthorized,
			Errors:  err.Error(),
		})
	}

	// create token
	token, err := helpers.CreateTokenJWT(user.ID)
	if err != nil {
		return helpers.InternalServerError(c, err)
	}

	// if user already exists just return user
	return c.JSON(helpers.SuccessResponse{
		Success: true,
		Message: "Signin success",
		Data: fiber.Map{
			"user":  user,
			"token": token,
		},
	})

}
