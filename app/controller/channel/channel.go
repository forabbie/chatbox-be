package controller

import (
	"chatbox/pkg/settings"
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"chatbox/pkg/util/validate"

	mchannel "chatbox/app/model/channel"
	schannel "chatbox/app/service/channel"

	jwtv4 "github.com/golang-jwt/jwt/v4"
)

func CreateChannel(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()

	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	claims, _ := c.Locals("claims").(jwtv4.MapClaims)

	sub, _ := claims["sub"].(float64)

	userId := int(sub)

	createdBy := int64(userId)

	payload := new(mchannel.CreatePayload)

	if err := c.BodyParser(payload); err != nil {
		log.Print(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if invalid := validate.All(payload); len(invalid) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"response": invalid})
	}

	channel, err := schannel.Insert(ctx, payload.Name, createdBy, payload.UserIDs)
	if err != nil {
		log.Print(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create channel")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"response": channel,
	})
}

func GetUserChannels(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()

	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	// Retrieve claims from context
	claimsValue := c.Locals("claims")
	claims, ok := claimsValue.(jwtv4.MapClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	// Extract user ID (subject) from claims
	sub, ok := claims["sub"].(float64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid subject in token")
	}

	userID := int64(sub)

	// 🔍 Query channels where user is a member
	channels, err := schannel.GetByUserID(ctx, userID)
	if err != nil {
		log.Println("Failed to get channels:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve channels")
	}

	return c.JSON(fiber.Map{
		"response": channels,
	})
}

func GetChannelDetailsByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
	defer cancel()

	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	channelIDParam := c.Params("id")
	channelID, err := strconv.ParseInt(channelIDParam, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid channel ID")
	}

	channel, err := schannel.GetDetailsByID(ctx, channelID)
	if err != nil {
		log.Println("Failed to get channel:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve channel")
	}

	return c.JSON(fiber.Map{
		"response": channel,
	})
}
