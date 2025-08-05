package routes

import (
	"dms-api/internal/database"
	"dms-api/internal/handlers"
	"dms-api/internal/repository"
	"dms-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(app fiber.Router) {

	dmsCountRepo := repository.DMSUserRepositorysInit(database.Database)
	dmsCountSvc := services.DMSUserServicesInit(dmsCountRepo)
	dmsCountHandler := handlers.DMSUsersHandlersInit(dmsCountSvc)

	app.Get("/dms/dashboard/dms-users", dmsCountHandler.GetAllDMSUsersCountHandler)
}
