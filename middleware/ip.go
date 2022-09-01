package middleware

import (
	"log"

	"github.com/dixiedream/authenticator/db"
	"github.com/dixiedream/authenticator/model"
	"github.com/gofiber/fiber/v2"
)

func GetIP(c *fiber.Ctx) string {
	headers := c.GetReqHeaders()
	ip := headers["X-Real-Ip"]
	if ip == "" {
		log.Println("Unable to obtain IP from headers")
		ip = c.IP()
	}

	return ip
}

func ValidateIP(c *fiber.Ctx) error {
	ip := GetIP(c)
	server := &model.Server{IpAddress: ip}
	if err := db.DB.First(&server).Error; err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusForbidden)
	}

    c.Locals("ip", ip)

	return c.Next()
}
