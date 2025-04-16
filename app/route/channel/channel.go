package route

import (
	"github.com/gofiber/fiber/v2"

	cchannel "chatbox/app/controller/channel"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Post("/channel", hjwt.ValidateAccessToken, cchannel.CreateChannel)
	router.Get("/channels", hjwt.ValidateAccessToken, cchannel.GetUserChannels)
	router.Get("/channel/:id", hjwt.ValidateAccessToken, cchannel.GetChannelDetailsByID)
	router.Post("/channel/add_member", hjwt.ValidateAccessToken, cchannel.AddMemberToChannel)
	router.Delete("/channel/:id", hjwt.ValidateAccessToken, cchannel.DeleteChannel)
	router.Put("/channel/leave", hjwt.ValidateAccessToken, cchannel.LeaveChannel)
}
