package route

import (
	"github.com/gofiber/fiber/v2"

	cchannel "chatbox/app/controller/channel"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Post("/channel", hjwt.ValidateRefreshToken, cchannel.CreateChannel)
	router.Get("/channel", hjwt.ValidateRefreshToken, cchannel.GetUserChannels)
}
