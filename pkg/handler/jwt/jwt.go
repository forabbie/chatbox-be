package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"chatbox/pkg/jwt"
	"chatbox/pkg/settings"
)

func ValidateAccessToken(c *fiber.Ctx) error {
	auth := jwt.Auth(c.Get(fiber.HeaderAuthorization), jwt.AuthScheme)

	if _, err := jwt.ParseToken(auth, jwt.AccessTokenKey); err != nil {
		log.Print(err)

		c.Set(fiber.HeaderWWWAuthenticate, jwt.AuthScheme)

		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}

func ValidateRefreshToken(c *fiber.Ctx) error {
	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	auth := jwt.Auth(c.Get(fiber.HeaderAuthorization), jwt.AuthScheme)

	claims, err := jwt.ParseToken(auth, jwt.RefreshTokenKey)
	if err != nil {
		log.Print(err)

		c.Set(fiber.HeaderWWWAuthenticate, jwt.AuthScheme)

		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Set claims to locals
	c.Locals("claims", claims)

	return c.Next()
}
