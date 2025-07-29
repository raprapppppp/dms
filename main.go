package main

import (
	"dms-api/database"
	"dms-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	database.ConnectionDB()
	app := fiber.New()
	
	routes.LoginRoute(app)

	app.Listen(":4000")

}