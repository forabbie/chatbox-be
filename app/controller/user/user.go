package controller

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
	"golang.org/x/crypto/bcrypt"

	"chatbox/pkg/jwt"
	"chatbox/pkg/settings"
	"chatbox/pkg/util/validate"

	muser "chatbox/app/model/user"
	suser "chatbox/app/service/user"
)

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
	defer cancel()

	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	emailaddress := c.FormValue("emailaddress")
	password := c.FormValue("password")

	// Validate input
	var validationErrors []validate.Map
	if v := validate.One("emailaddress", emailaddress, "required,emailaddress"); len(v) > 0 {
		validationErrors = append(validationErrors, v)
	}
	if v := validate.One("password", password, "required"); len(v) > 0 {
		validationErrors = append(validationErrors, v)
	}
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"response": validationErrors,
		})
	}

	// Check for existing user
	filter := map[string][]string{
		"or": {"emailaddress = $1"},
	}
	args := []interface{}{emailaddress}

	count, err := suser.Count(ctx, filter, args)
	if err != nil {
		log.Print(err)
		return err
	}
	if count > 0 {
		return c.SendStatus(fiber.StatusConflict)
	}

	// Create new user
	user := &muser.User{
		Firstname:    c.FormValue("firstname"),
		Lastname:     c.FormValue("lastname"),
		Username:     c.FormValue("username"),
		EmailAddress: emailaddress,
	}
	isActive := c.FormValue("is_active") == "true"
	user.IsActive = &isActive

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return err
	}
	user.Password = string(hashedPassword)

	// Insert user
	userID, err := suser.Insert(ctx, user)
	if err != nil {
		log.Print(err)
		return err
	}

	// Generate tokens
	accessToken, err := jwt.NewToken(userID, settings.ShortExpiration, nil, jwt.AccessTokenKey)
	if err != nil {
		log.Print(err)
		return err
	}

	refreshToken, err := jwt.NewToken(userID, settings.LongExpiration, utils.UUID(), jwt.RefreshTokenKey)
	if err != nil {
		log.Print(err)
		return err
	}

	// Respond with tokens
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"response": fiber.Map{
			"access_token":  accessToken,
			"token_type":    jwt.AuthScheme,
			"expires_in":    settings.ShortExpiration.Seconds(),
			"refresh_token": refreshToken,
		},
	})
}
