package routes

import (
	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(app fiber.Router){


	app.Get("/dms/dashboard/get-all-dms-user")
}