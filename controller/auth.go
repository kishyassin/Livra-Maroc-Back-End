package controller

import (
	"kishyassin/Livra-Maroc/model"
	"kishyassin/Livra-Maroc/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login handles user authentication via email and password
func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Check if user exists
		var user model.User
		if err := db.Where("email = ?", request.Email).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		// Generate JWT tokens
		accessToken, accessExpiresAt, err := utils.GenerateJWT(user.ID, user.Role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not login",
			})
		}

		refreshToken, _, err := utils.GenerateRefreshToken(user.ID, user.Role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not login",
			})
		}

		// Return response
		return c.JSON(fiber.Map{
			"access_token":      accessToken,
			"refresh_token":     refreshToken,
			"access_expires_at": accessExpiresAt,
		})
	}
}
