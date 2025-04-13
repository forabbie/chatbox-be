package controller

import (
	"context"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
	"golang.org/x/crypto/bcrypt"

	jwtv4 "github.com/golang-jwt/jwt/v4"

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

func Login(c *fiber.Ctx) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)
	defer cancel()

	// Prevent caching
	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	// Extract form values
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Validate input
	var validationErrors []validate.Map
	if err := validate.One("username", username, "required"); len(err) > 0 {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.One("password", password, "required"); len(err) > 0 {
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"response": validationErrors})
	}

	// Build query filter
	filter := map[string][]string{
		"or": {
			"emailaddress = ?",
			"username = ?",
		},
	}
	args := []interface{}{username, username}

	// Fetch account
	limit := 1
	users, err := suser.Fetch(ctx, filter, args, limit)
	if err != nil {
		log.Printf("failed to fetch account: %v", err)
		return fiber.ErrInternalServerError
	}
	if len(users) == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	user := users[0]

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Check user activation
	if user.IsActive != nil && !*user.IsActive {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// Generate access token
	accessToken, err := jwt.NewToken(
		user.Id,
		settings.ShortExpiration,
		nil,
		jwt.AccessTokenKey,
	)
	if err != nil {
		log.Printf("failed to generate access token: %v", err)
		return fiber.ErrInternalServerError
	}

	// Generate refresh token
	refreshToken, err := jwt.NewToken(
		user.Id,
		settings.LongExpiration,
		utils.UUID(),
		jwt.RefreshTokenKey,
	)
	if err != nil {
		log.Printf("failed to generate refresh token: %v", err)
		return fiber.ErrInternalServerError
	}

	// Success response
	return c.JSON(fiber.Map{
		"response": fiber.Map{
			"access_token":  accessToken,
			"token_type":    jwt.AuthScheme,
			"expires_in":    settings.ShortExpiration.Seconds(),
			"refresh_token": refreshToken,
		},
	})
}

func Refresh(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()

	claims := c.Locals("claims").(jwtv4.MapClaims)

	sub, _ := claims["sub"].(float64)

	id := int(sub)

	user, err := suser.Get(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)

			return c.SendStatus(fiber.StatusNotFound)
		} else {
			log.Print(err)

			return err
		}
	}

	if user.IsActive != nil && !*user.IsActive {
		return c.SendStatus(fiber.StatusForbidden)
	}

	accessToken, err := jwt.NewToken(
		user.Id,
		settings.ShortExpiration,
		nil,
		jwt.AccessTokenKey,
	)
	if err != nil {
		log.Print(err)

		return err
	}

	refreshToken, err := jwt.NewToken(
		user.Id,
		settings.LongExpiration,
		utils.UUID(),
		jwt.RefreshTokenKey,
	)
	if err != nil {
		log.Print(err)

		return err
	}

	return c.JSON(fiber.Map{
		"response": fiber.Map{
			"access_token":  accessToken,
			"token_type":    jwt.AuthScheme,
			"expires_in":    settings.ShortExpiration.Seconds(),
			"refresh_token": refreshToken,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
