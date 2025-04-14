package controller

import (
	"chatbox/pkg/settings"
	"chatbox/pkg/util/validate"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

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

func GetMessages(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
	defer cancel()

	query := new(mdmsg.Query)

	if err := c.QueryParser(query); err != nil {
		log.Print(err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	filter := map[string][]string{
		"or":  {},
		"and": {},
	}
	args := []interface{}{}

	// Fulltext search via `q`
	if q := strings.TrimSpace(c.Query("q")); q != "" {
		for _, field := range []string{"message", "sender.firstname", "sender.lastname", "sender.username"} {
			filter["or"] = append(filter["or"], fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, "%"+q+"%")
		}
	}

	// Field-based filters
	// if query.Message != nil {
	// 	filter["or"] = append(filter["or"], "dm.message ILIKE ?")
	// 	args = append(args, "%"+*query.Message+"%")
	// }
	// if query.Firstname != nil {
	// 	filter["or"] = append(filter["or"], "sender.firstname ILIKE ?")
	// 	args = append(args, "%"+*query.Firstname+"%")
	// }
	// if query.Lastname != nil {
	// 	filter["or"] = append(filter["or"], "sender.lastname ILIKE ?")
	// 	args = append(args, "%"+*query.Lastname+"%")
	// }
	// if query.Username != nil {
	// 	filter["or"] = append(filter["or"], "sender.username ILIKE ?")
	// 	args = append(args, "%"+*query.Username+"%")
	// }
	if query.SenderID != nil && query.ReceiverID != nil {
		filter["or"] = append(filter["or"],
			"(dm.sender_id = ? AND dm.receiver_id = ?)",
			"(dm.sender_id = ? AND dm.receiver_id = ?)",
		)
		args = append(args, *query.SenderID, *query.ReceiverID, *query.ReceiverID, *query.SenderID)
	}
	if query.Created.Gte != nil {
		filter["and"] = append(filter["and"], "dm.sent_at >= ?")
		args = append(args, *query.Created.Gte)
	}
	if query.Created.Lte != nil {
		filter["and"] = append(filter["and"], "dm.sent_at <= ?")
		args = append(args, *query.Created.Lte)
	}

	if len(filter["or"]) == 0 {
		delete(filter, "or")
	}
	if len(filter["and"]) == 0 {
		delete(filter, "and")
	}

	// Sort and pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	sorts := strings.Split(c.Query("sort"), ",")
	order, sort := "dm.sent_at", "ASC"

	if len(sorts) >= 1 && strings.TrimSpace(sorts[0]) != "" {
		// Example validation: only allow specific fields
		switch sorts[0] {
		case "sent_at", "firstname", "lastname":
			order = "dm." + sorts[0]
		}
	}
	if len(sorts) == 2 {
		if strings.ToLower(sorts[1]) == "desc" {
			sort = "DESC"
		}
	}

	// Fetch results
	messages, err := sdmsg.Fetch(ctx, filter, args, order, sort, limit, offset)
	if err != nil {
		log.Print(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch messages")
	}

	return c.JSON(fiber.Map{
		"response": messages,
		"total":    len(messages), // optional: you can create a Count method like in your user example
	})
}
