package middleware

import (
	"dms-api/config"
	"dms-api/utils/cryptography/security"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func PasswordResetAuthMiddleware(m *fiber.Ctx) error {
	//Get Cookies
	cookie := m.Cookies("otpToken")
	if cookie == "" {
		return m.Status(fiber.StatusUnauthorized).SendString("No token cookie found")
	}
	//Create Instance of NewAPISecurity
	apiSec, _ := security.NewAPISecurity()
	//Decrypt the token
	decryptedToken, _ := apiSec.Decrypt(cookie)
	//Convert to string
	convertedToken := string(decryptedToken)
	//Check and Validate
	token, err := jwt.Parse(convertedToken, func(t *jwt.Token) (interface{}, error) {

		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return []byte(config.Config("COOKIES_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return m.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}
	//Check the puspose of the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// You can check the "purpose" claim here if needed
		if claims["purpose"] != "password_reset" {
			return m.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid token purpose",
			})
		}
		//Get the id of the requestor from claims and store to locals
		id := claims["user_id"]
		m.Locals("user_id", id)
	}
	return m.Next()
}
