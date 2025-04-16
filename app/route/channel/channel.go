package route

import (
	"github.com/gofiber/fiber/v2"

	cchannel "chatbox/app/controller/channel"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Post("/channels", hjwt.ValidateRefreshToken, cchannel.CreateChannel)
	router.Get("/channels", hjwt.ValidateRefreshToken, cchannel.GetUserChannels)
	router.Get("/channels/:id", hjwt.ValidateRefreshToken, cchannel.GetChannelDetailsByID)
	router.Post("/channels/add_member", hjwt.ValidateRefreshToken, cchannel.AddMemberToChannel)
}
