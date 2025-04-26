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

	mmsg "chatbox/app/model/message"
	smsg "chatbox/app/service/message"

	jwtv4 "github.com/golang-jwt/jwt/v4"
)

func SendMessage(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
	defer cancel()
	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	// Get sender from JWT
	claims, _ := c.Locals("claims").(jwtv4.MapClaims)
	sub, _ := claims["sub"].(float64)
	senderID := int64(sub)

	// Parse the request body
	msg := new(mmsg.Message)
	if err := c.BodyParser(msg); err != nil {
		log.Print(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	// Set the authenticated sender ID
	msg.Sender.ID = senderID

	// Validate required fields (receiver_id, receiver_class, message)
	if invalid := validate.All(msg); len(invalid) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"response": invalid})
	}

	// Ensure receiver_class is either "user" or "channel"
	msg.ReceiverClass = strings.ToLower(msg.ReceiverClass)
	if msg.ReceiverClass != "user" && msg.ReceiverClass != "channel" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid receiver_class. Must be 'user' or 'channel'.",
		})
	}

	// Insert the message
	result, err := smsg.Insert(ctx, msg)
	if err != nil {
		log.Print(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send message")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"response": fiber.Map{
			"id":      result.ID,
			"sent_at": result.SentAt,
		},
	})
}

func GetMessages(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
	defer cancel()

	claims, _ := c.Locals("claims").(jwtv4.MapClaims)
	sub, _ := claims["sub"].(float64)
	userId := int(sub)
	requestBy := int64(userId)

	query := new(mmsg.Query)
	if err := c.QueryParser(query); err != nil {
		log.Print(err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	filter := map[string][]string{
		"or":  {},
		"and": {},
	}
	args := []interface{}{}

	// Fulltext search
	if q := strings.TrimSpace(c.Query("q")); q != "" {
		for _, field := range []string{"message", "sender.firstname", "sender.lastname", "sender.username"} {
			filter["or"] = append(filter["or"], fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, "%"+q+"%")
		}
	}

	// Validate receiver_class only if present
	if query.ReceiverClass != nil && *query.ReceiverClass != "user" && *query.ReceiverClass != "channel" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid receiver_class. Must be 'user' or 'channel'.",
		})
	}

	// Match conversation: user ↔ receiver
	if query.ReceiverID != nil && query.ReceiverClass != nil {
		filter["or"] = append(filter["or"],
			"((dm.sender_id = ? AND dm.receiver_id = ?) OR (dm.sender_id = ? AND dm.receiver_id = ?))",
		)
		args = append(args, requestBy, *query.ReceiverID, *query.ReceiverID, requestBy)
	}

	// Date filters
	if query.Created.Gte != nil {
		filter["and"] = append(filter["and"], "dm.sent_at >= ?")
		args = append(args, *query.Created.Gte)
	}
	if query.Created.Lte != nil {
		filter["and"] = append(filter["and"], "dm.sent_at <= ?")
		args = append(args, *query.Created.Lte)
	}

	// Clean up empty filters
	if len(filter["or"]) == 0 {
		delete(filter, "or")
	}
	if len(filter["and"]) == 0 {
		delete(filter, "and")
	}

	// Pagination and sorting
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	sorts := strings.Split(c.Query("sort"), ",")
	order, sort := "dm.sent_at", "ASC"
	if len(sorts) >= 1 && strings.TrimSpace(sorts[0]) != "" {
		switch sorts[0] {
		case "sent_at", "firstname", "lastname":
			order = "dm." + sorts[0]
		}
	}
	if len(sorts) == 2 && strings.ToLower(sorts[1]) == "desc" {
		sort = "DESC"
	}

	// Fetch results
	messages, err := smsg.Fetch(ctx, filter, args, order, sort, limit, offset)
	if err != nil {
		log.Print(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch messages")
	}

	return c.JSON(fiber.Map{
		"response": messages,
		"total":    len(messages),
	})
}
