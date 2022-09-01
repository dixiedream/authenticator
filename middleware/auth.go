package middleware

import (
	"os"

	"github.com/dixiedream/authenticator/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// Protect protected routes
func Auth() fiber.Handler {
	signingKey := os.Getenv("JWT_ACCESS_SECRET")
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(signingKey),
		Claims:       utils.Claim{},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "Missing or malformed JWT"})
	} else {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "Invalid or expired JWT"})
	}
}
