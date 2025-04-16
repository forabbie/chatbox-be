package route

import (
	"github.com/gofiber/fiber/v2"

	cchannel "chatbox/app/controller/channel"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Post("/channels", hjwt.ValidateAccessToken, cchannel.CreateChannel)
	router.Get("/channels", hjwt.ValidateAccessToken, cchannel.GetUserChannels)
	router.Get("/channels/:id", hjwt.ValidateAccessToken, cchannel.GetChannelDetailsByID)
	router.Post("/channels/add_member", hjwt.ValidateAccessToken, cchannel.AddMemberToChannel)
}
