package handler

import (
	"log"

	"github.com/dixiedream/authenticator/db"
	"github.com/dixiedream/authenticator/model"
	"github.com/gofiber/fiber/v2"
)

func CreateServer(c *fiber.Ctx) error {
	var input model.Server
	if err := c.BodyParser(&input); err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}

	name := input.Name
	hostname := input.Hostname
	password := input.Password
	ipAddress := input.IpAddress

	if name == "" || hostname == "" || password == "" || ipAddress == "" {
		log.Println("Missing values")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := db.DB.Create(&input).Error; err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        input.ID,
		"hostname":  input.Hostname,
		"ipAddress": input.IpAddress,
		"role":      input.Role,
	})
}
