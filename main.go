package main

import (
	"log"

	"github.com/dixiedream/authenticator/db"
	"github.com/dixiedream/authenticator/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	db.Connect()

	router.Setup(app)
	log.Fatal(app.Listen(":3000"))
}
