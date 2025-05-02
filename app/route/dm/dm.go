package route

import (
	"github.com/gofiber/fiber/v2"

	cdm "chatbox/app/controller/dm"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Get("/direct-messages", hjwt.ValidateAccessToken, cdm.GetUserDMList)
}
