package route

import (
	"github.com/gofiber/fiber/v2"

	cdmsg "chatbox/app/controller/message"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Post("/message", hjwt.ValidateRefreshToken, cdmsg.SendMessage)
}
