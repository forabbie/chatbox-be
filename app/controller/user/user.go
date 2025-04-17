package controller

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
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
	accessToken, err := jwt.NewToken(userID, settings.ShortExpiration, nil, os.Getenv("JWT_ACCESS_TOKEN_KEY"))
	if err != nil {
		log.Print(err)
		return err
	}

	refreshToken, err := jwt.NewToken(userID, settings.LongExpiration, utils.UUIDv4(), os.Getenv("JWT_REFRESH_TOKEN_KEY"))
	if err != nil {
		log.Print(err)
		return err
	}

	// Respond with tokens
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"response": fiber.Map{
			"access_token":  accessToken,
			"token_type":    settings.BearerAuthScheme,
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
		log.Printf("failed to fetch user: %v", err)
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
		os.Getenv("JWT_ACCESS_TOKEN_KEY"),
	)
	if err != nil {
		log.Printf("failed to generate access token: %v", err)
		return fiber.ErrInternalServerError
	}

	refreshTokenExpiration := settings.LongExpiration

	// Generate refresh token
	refreshToken, err := jwt.NewToken(
		user.Id,
		refreshTokenExpiration,
		utils.UUIDv4(),
		os.Getenv("JWT_REFRESH_TOKEN_KEY"),
	)
	if err != nil {
		log.Print(err)

		return err
	}

	return c.JSON(fiber.Map{
		"response": fiber.Map{
			"access_token":  accessToken,
			"token_type":    settings.BearerAuthScheme,
			"expires_in":    settings.ShortExpiration.Seconds(),
			"refresh_token": refreshToken,
		},
	})
}

func Refresh(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()

	claims := c.Locals("claims").(jwtv4.MapClaims)
	userID := int64(claims["sub"].(float64))

	user, err := suser.GetByID(ctx, userID)
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
		os.Getenv("JWT_ACCESS_TOKEN_KEY"),
	)
	if err != nil {
		log.Print(err)

		return err
	}

	auth := jwt.ParseAuth(c.Get(fiber.HeaderAuthorization), settings.BearerAuthScheme)

	refreshToken := auth

	return c.JSON(fiber.Map{
		"response": fiber.Map{
			"access_token":  accessToken,
			"token_type":    settings.BearerAuthScheme,
			"expires_in":    settings.ShortExpiration.Seconds(),
			"refresh_token": refreshToken,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()

	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	users, err := suser.GetAll(ctx)
	if err != nil {
		log.Println("Failed to retrieve users:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve users")
	}

	return c.JSON(fiber.Map{
		"response": users,
	})
}

func GetUserDetails(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()

	c.Set(fiber.HeaderCacheControl, settings.CacheControlNoStore)

	userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := suser.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve user")
	}

	return c.JSON(fiber.Map{
		"response": user,
	})
}

func GetCurrentUser(c *fiber.Ctx) error {
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

	user, err := suser.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve user")
	}

	return c.JSON(fiber.Map{
		"response": user,
	})
}
