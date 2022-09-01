package handler

import (
	"log"
	"time"

	"github.com/dixiedream/authenticator/db"
	// "github.com/dixiedream/authenticator/middleware"
	"github.com/dixiedream/authenticator/model"
	"github.com/dixiedream/authenticator/utils"
	"github.com/gofiber/fiber/v2"
)

const REFRESH_TOKEN_COOKIE_NAME = "refresh_token"

func Refresh(c *fiber.Ctx) error {
	cookie := c.Cookies(REFRESH_TOKEN_COOKIE_NAME)
	_, err := utils.RefreshTokenIsValid(cookie)
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}
	session := model.Session{}
	if err := db.DB.Where(&model.Session{Token: cookie, Expired: false}).First(&session).Error; err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	server := model.Server{}
	if err := db.DB.First(&server, session.ServerID).Error; err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	token, err := utils.GenerateAccessToken(&utils.Payload{Hostname: server.Hostname, Role: server.Role})
	if err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": token})
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ip := c.Locals("ip")
	pass := input.Password

	server := model.Server{}
	if err := db.DB.Where(&model.Server{IpAddress: ip.(string)}).First(&server).Error; err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if err := utils.PasswordIsValid(server.Password, pass); err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	db.DB.Model(&server).Updates(model.Server{LastAccess: time.Now()})

	token, err := utils.GenerateAccessToken(&utils.Payload{Hostname: server.Hostname, Role: server.Role})
	if err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	rToken, err := utils.GenerateRefreshToken(&utils.Payload{Hostname: server.Hostname})
	if err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	session := model.Session{Token: rToken, ServerID: server.ID, Expired: false}
	if err := db.DB.Create(&session).Error; err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := &fiber.Cookie{
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		Name:     REFRESH_TOKEN_COOKIE_NAME,
		Value:    rToken,
		HTTPOnly: true,
	}
	c.Cookie(cookie)
	return c.JSON(fiber.Map{"token": token})
}

func Logout(c *fiber.Ctx) error {
	rToken := c.Cookies(REFRESH_TOKEN_COOKIE_NAME)
	if rToken == "" {
		return c.SendStatus(fiber.StatusNoContent)
	}

	session := model.Session{Token: rToken}

	db.DB.Model(&session).Updates(model.Session{Expired: true})
	if err := db.DB.Delete(&session).Error; err != nil {
		log.Println(err.Error())
	}

	c.ClearCookie(REFRESH_TOKEN_COOKIE_NAME)
	return c.SendStatus(fiber.StatusNoContent)
}
