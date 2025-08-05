package routes

import (
	"dms-api/internal/database"
	"dms-api/internal/handlers"
	"dms-api/internal/repository"
	"dms-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

func DocumentTypesRoutes(app fiber.Router) {

	//Initialize repo,service,handler
	dtRepo := repository.DocumentTypesRepositoryInit(database.Database)
	dtServices := services.DocumentTypesServicesInit(dtRepo)
	dtHandler := handlers.DocumentTypesHandlersInit(dtServices)

	app.Get("/dms/get-document-types", dtHandler.GetAllDocumentTypesHandler)

}