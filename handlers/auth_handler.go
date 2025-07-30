package handlers

import (
	"dms-api/modals"
	"dms-api/services"
	"dms-api/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type InjectLoginHandler struct {
	services services.LoginServices
}

func LoginHandlerInit(s services.LoginServices) *InjectLoginHandler {
	return &InjectLoginHandler{s}
}
var parseError = "Cannot parse JSON"
// Login
func (h *InjectLoginHandler) LoginHandler(hh *fiber.Ctx) error {
	var cred modals.Login
	if err := hh.BodyParser(&cred); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": parseError})
	}
	user, err := h.services.LoginService(cred)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": utils.ErrNotFound})
		}
		if errors.Is(err, utils.ErrNotMatch) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": utils.ErrNotMatch})
		}
	}
	return hh.Status(fiber.StatusOK).JSON(fiber.Map{"message": user})
}

// Register
func (h *InjectLoginHandler) Registerhandler(hh *fiber.Ctx) error {
	var cred modals.Accounts
	if err := hh.BodyParser(&cred); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": parseError})
	}
	_, err := h.services.RegisterService(cred)
	if err != nil {
		if errors.Is(err, utils.ErrAlreadyExist) {
			return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": utils.ErrAlreadyExist})
		}
	}
	return hh.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Created"})
}

// Forgot Password
func (h *InjectLoginHandler) ForgotPasswordRequestHandler(hh *fiber.Ctx) error {
	var email modals.Forgot

	if err := hh.BodyParser(&email); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": parseError})
	}

	mess, err := h.services.ForgotPasswordRequestService(email)
	if err != nil {
		if errors.Is(err, utils.ErrEmailNotExist) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": utils.ErrEmailNotExist})
		}
		if errors.Is(err, utils.ErrOTPGenerationFailed) {
			return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": utils.ErrOTPGenerationFailed})
		}
		if errors.Is(err, utils.ErrSendingOTP) {
			return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": utils.ErrSendingOTP})
		}
		if errors.Is(err, utils.ErrOTPRequestLimit) {
			return hh.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"message": utils.ErrOTPRequestLimit})
		}
	}
	return hh.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": mess})
}
//Verify OTP
func(h *InjectLoginHandler) VerifyOTPHandler(hh *fiber.Ctx) error {
	var otp modals.VerifyOTP

	if err := hh.BodyParser(&otp); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": parseError}) 
	}
	mess, err := h.services.VerifyOTPService(otp)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound){
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Email not found"})
		}
		if errors.Is(err, utils.ErrNoLatestOtp){
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No Password reset request"})
		}
		if errors.Is(err, utils.ErrOTPAlreadyUsed){
			return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "OTP already used"})
		}
		if errors.Is(err, utils.ErrOTPExpired){
			return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "OTP Expired"})
		}
		if errors.Is(err, utils.ErrNotMatch){
			return hh.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "OTP not match"})
		}
		if errors.Is(err, utils.ErrOTPUpdateFailed){
			return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "failed to updated OTP status"})
		}
	}
	return hh.Status(fiber.StatusOK).JSON(fiber.Map{"message": mess})
}
