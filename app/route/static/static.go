package route

import (
	"github.com/gofiber/fiber/v2"

	hfilesystem "chatbox/pkg/handler/filesystem"
	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Use("/static", hjwt.ValidateAccessToken, hfilesystem.Static)
}
