package service

import (
	"chatbox/pkg/settings"
	"context"
	"log"

	sdm "chatbox/app/service/dm"

	"github.com/gofiber/fiber/v2"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

func GetUserDMList(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
	defer cancel()

	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	claimsValue := c.Locals("claims")
	claims, ok := claimsValue.(jwtv4.MapClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid subject in token")
	}

	userID := int64(sub)

	dms, err := sdm.GetDMListByUserID(ctx, userID)
	if err != nil {
		log.Println("Failed to get DM list:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve DMs")
	}

	return c.JSON(fiber.Map{
		"response": dms,
	})
}
