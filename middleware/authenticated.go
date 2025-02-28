package middleware

import (
	"kishyassin/Livra-Maroc/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// authenticated middleware, validate JWT token
func Authenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header missing",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		accessToken := parts[1]

		_, claims, err := utils.ValidateToken(accessToken, false)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired access token",
			})
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}

		c.Locals("userID", uint(userID))
		return c.Next()
	}
}
