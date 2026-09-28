package validation

import (
	"github.com/RomeyKung/social/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var ValidateUser = validator.New()

func validateRequestBody(c *fiber.Ctx, body interface{}) error {
	var errors []*models.IError

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	err := ValidateUser.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.IError
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	return c.Next()
}

func ValidateCreateUser(c *fiber.Ctx) error {
	var body models.CreateUser
	return validateRequestBody(c, &body)
}

func ValidateLoginUser(c *fiber.Ctx) error {
	var body models.LoginUser
	return validateRequestBody(c, &body)
}
