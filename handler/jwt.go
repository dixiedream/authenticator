package handler

import (
	"log"

	"github.com/dixiedream/authenticator/db"
	"github.com/dixiedream/authenticator/model"
	"github.com/dixiedream/authenticator/utils"
	"github.com/gofiber/fiber/v2"
)

func CheckToken(c *fiber.Ctx) error {
    type CheckTokenInput struct {
        Token   string  `json:"token"`
    }
    input := &CheckTokenInput{}
	if err := c.BodyParser(&input); err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}
    
    payload, err := utils.AccessTokenIsValid(input.Token)
    if err != nil {
        log.Println(err.Error())
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    server := &model.Server{Hostname: payload.Hostname}
    if err := db.DB.First(&server).Error; err != nil {
        log.Println(err.Error())
        return c.SendStatus(fiber.StatusUnauthorized)
    }
        
    return c.SendStatus(fiber.StatusAccepted)
}
