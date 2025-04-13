package route

import (
	"github.com/gofiber/fiber/v2"

	hfile "chatbox/pkg/handler/file"
	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Post("/file", hjwt.ValidateAccessToken, hfile.Upload)

	router.Get("/file/:filename", hjwt.ValidateAccessToken, hfile.Download)
}
