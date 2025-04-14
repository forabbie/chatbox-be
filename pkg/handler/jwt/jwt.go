package handler

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"chatbox/pkg/jwt"
	"chatbox/pkg/settings"
)

func ValidateAccessToken(c *fiber.Ctx) error {
	accessToken := jwt.ParseAuth(c.Get(fiber.HeaderAuthorization), settings.BearerAuthScheme)

	accessToken = c.Query("token", accessToken)

	claims, err := jwt.ParseToken(accessToken, os.Getenv("JWT_ACCESS_TOKEN_KEY"))
	if err != nil {
		log.Print(err)

		c.Set(fiber.HeaderWWWAuthenticate, settings.BearerAuthScheme)

		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("access_token", accessToken)

	// Set claims to locals
	c.Locals("claims", claims)

	return c.Next()
}

func ValidateRefreshToken(c *fiber.Ctx) error {
	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	refreshToken := jwt.ParseAuth(c.Get(fiber.HeaderAuthorization), settings.BearerAuthScheme)

	claims, err := jwt.ParseToken(refreshToken, os.Getenv("JWT_REFRESH_TOKEN_KEY"))
	if err != nil {
		log.Print(err)

		c.Set(fiber.HeaderWWWAuthenticate, settings.BearerAuthScheme)

		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("refresh_token", refreshToken)

	// Set claims to locals
	c.Locals("claims", claims)

	return c.Next()
}
