package handlers

import (
	"dms-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type InjectDMSUsersServices struct {
	service services.DMSUsersServices
}

// Handler Initializer
func DMSUsersHandlersInit(s services.DMSUsersServices) *InjectDMSUsersServices {
	return &InjectDMSUsersServices{s}
}

func (h InjectDMSUsersServices) GetAllDMSUsersCountHandler(hh *fiber.Ctx) error {
	return hh.SendStatus(fiber.StatusOK)
}
