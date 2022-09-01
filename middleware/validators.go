package middleware

import (
	"github.com/dixiedream/authenticator/model"
	"github.com/dixiedream/authenticator/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validator = validator.New()

func ValidateCreateServer(c *fiber.Ctx) error {
	var errors []*utils.IError
	body := new(model.Server)
	c.BodyParser(&body)

	if err := Validator.Struct(body); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			el := utils.IError{Field: err.Field(), Tag: err.Tag(), Value: err.Param()}
			errors = append(errors, &el)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	return c.Next()
}
