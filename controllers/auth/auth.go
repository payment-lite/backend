package authController

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
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
	//result := database.DBConn.Model(&models.User{}).Create(&models.User{
	//	Name:     reqData.Name,
	//	Email:    reqData.Email,
	//	Password: string(hashedPassword),
	//}).First(&user)

	result := database.DBConn.Transaction(func(tx *gorm.DB) error {
		//	create user
		if err := tx.Model(&models.User{}).Create(&models.User{
			Name:     reqData.Name,
			Email:    reqData.Email,
			Password: string(hashedPassword),
		}).First(&user).Error; err != nil {
			return err
		}

		//	create team
		team := &models.Team{
			Name:    reqData.Name,
			OwnerID: &user.ID,
		}
		if err := tx.Model(&models.Team{}).Create(&team).Error; err != nil {
			return err
		}

		// Set User's TeamID
		user.TeamID = &team.ID
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if result != nil {
		fmt.Println(result.(*mysql.MySQLError))
		if result.(*mysql.MySQLError).Number == 1062 {
			return c.Status(fiber.StatusBadRequest).JSON(&helpers.ErrorResponse{
				Success: false,
				Message: "Email already exists",
				Errors:  nil,
			})
		}
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
// besok terusin gimana cara nextauth consume token yg kita berikan
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

	var user models.User

	// find user by email
	result := database.DBConn.Where("email = ?", reqData.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			var user models.User
			user.Name = reqData.Name
			user.Email = reqData.Email
			user.Password = string(hashedPassword)
			result := database.DBConn.Model(&models.User{}).Create(&user)
			if result.Error != nil {
				return helpers.InternalServerError(c, result.Error)
			}

			//create team
			result = database.DBConn.Model(&models.Team{}).Create(&models.Team{
				Name:    user.Name,
				OwnerID: &user.ID,
			})
			if result.Error != nil {
				return helpers.InternalServerError(c, result.Error)
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
		return helpers.InternalServerError(c, result.Error)
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
