package controller

import (
	"chatbox/pkg/settings"
	"context"
	"log"

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
