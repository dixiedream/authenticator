package middleware

import (
	"log"

	"github.com/dixiedream/authenticator/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func IsAdmin(c *fiber.Ctx) error {
	log.Println(c.Locals("user"))
	token := c.Locals("user").(*jwt.Token)
    if claims, ok := token.Claims.(utils.Claim); ok && token.Valid {
        log.Println(&utils.Payload{Hostname: claims.Name, Role: claims.Role})
	} else {
        return c.SendStatus(401)
	}

	return c.Next()
}
