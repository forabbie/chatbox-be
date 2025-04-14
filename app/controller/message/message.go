package controller

import (
	"chatbox/pkg/settings"
	"chatbox/pkg/util/validate"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	mdmsg "chatbox/app/model/message"
	sdmsg "chatbox/app/service/message"
)

func SendMessage(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()
	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	dmsg := new(mdmsg.DirectMessage)

	if err := c.BodyParser(dmsg); err != nil {
		log.Print(err)

		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if invalid := validate.All(dmsg); len(invalid) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"response": invalid})
	}

	lastInsertId, err := sdmsg.Insert(ctx, dmsg)
	if err != nil {
		log.Print(err)

		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"response": fiber.Map{"lastInsertId": lastInsertId}})
}
