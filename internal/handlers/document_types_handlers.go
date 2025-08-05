package handlers

import (
	"dms-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type InjectDocumentTypesServices struct {
	service services.DocumentTypesServices
}

// Handler Initializer
func DocumentTypesHandlersInit(s services.DocumentTypesServices) *InjectDocumentTypesServices {
	return &InjectDocumentTypesServices{s}
}

// Get all Document types
func (h *InjectDocumentTypesServices) GetAllDocumentTypesHandler(dt *fiber.Ctx) error {
	allDocumentTypes, err := h.service.GetAllDocumentTypes()
	if err != nil {
		return err
	}
	return dt.Status(fiber.StatusOK).JSON(allDocumentTypes)
}
