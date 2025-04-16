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

func AddMemberToChannel(c *fiber.Ctx) error {
	var req mchannel.AddMemberRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.ID == 0 || req.MemberID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Missing id or member_id")
	}

	// Optional: check if the user is already a member
	exists, err := schannel.IsMember(req.ID, req.MemberID)
	if err != nil {
		log.Println("Error checking membership:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check membership")
	}
	if exists {
		return fiber.NewError(fiber.StatusConflict, "User is already a member of the channel")
	}

	// Insert new member into channel
	err = schannel.AddMember(req.ID, req.MemberID)
	if err != nil {
		log.Println("Failed to add member:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to add member to channel")
	}

	return c.JSON(fiber.Map{
		"message": "Member added successfully",
	})
}

func LeaveChannel(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
	defer cancel()

	type Payload struct {
		ID int64 `json:"id"` // channel ID
	}

	var payload Payload
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid payload")
	}

	claims := c.Locals("claims").(jwtv4.MapClaims)
	userID := int64(claims["sub"].(float64))

	err := schannel.RemoveMember(ctx, payload.ID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to leave channel")
	}

	return c.JSON(fiber.Map{"message": "Left the channel successfully"})
}

func DeleteChannel(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
	defer cancel()

	claims := c.Locals("claims").(jwtv4.MapClaims)
	userID := int64(claims["sub"].(float64))
	channelID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid channel ID")
	}

	// Check if user is creator
	isCreator, err := schannel.IsChannelCreator(ctx, channelID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check channel owner")
	}
	if !isCreator {
		return fiber.NewError(fiber.StatusForbidden, "Only the channel creator can delete the channel")
	}

	if err := schannel.Delete(ctx, channelID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete channel")
	}

	return c.JSON(fiber.Map{"message": "Channel deleted successfully"})
}
