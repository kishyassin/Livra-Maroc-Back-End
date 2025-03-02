package controller

import (
	"kishyassin/Livra-Maroc/model"
	"kishyassin/Livra-Maroc/utils"
	"strings"

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

// RefreshToken handles token refresh
func RefreshToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "En-tête d'autorisation manquant",
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Format d'autorisation invalide",
			})
		}

		refreshToken := parts[1]

		// Validate refresh token
		_, claims, err := utils.ValidateToken(refreshToken, true)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Jeton de rafraîchissement invalide ou expiré",
			})
		}

		// Extract user details
		userID := uint(claims["user_id"].(float64)) // Convert float64 to uint
		role := claims["role"].(string)

		// Generate new access token
		newAccessToken, exp, err := utils.GenerateJWT(userID, role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Échec de la génération d'un nouveau jeton d'accès",
			})
		}

		// Generate new refresh token
		newRefreshToken, _, err := utils.GenerateRefreshToken(userID, role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Échec de la génération d'un nouveau jeton de rafraîchissement",
			})
		}

		return c.JSON(fiber.Map{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
			"exp":           exp,
		})
	}
}

