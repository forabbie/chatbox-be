package route

import (
	"github.com/gofiber/fiber/v2"

	cdmsg "chatbox/app/controller/message"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Get("/message", hjwt.ValidateAccessToken, cdmsg.GetMessages)
	router.Post("/message", hjwt.ValidateAccessToken, cdmsg.SendMessage)
}
