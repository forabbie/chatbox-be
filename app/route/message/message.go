package route

import (
	"github.com/gofiber/fiber/v2"

	cmsg "chatbox/app/controller/message"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Get("/message", hjwt.ValidateAccessToken, cmsg.GetMessages)
	router.Post("/message", hjwt.ValidateAccessToken, cmsg.SendMessage)
}
