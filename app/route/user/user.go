package route

import (
	"github.com/gofiber/fiber/v2"

	cuser "chatbox/app/controller/user"
)

func Route(router fiber.Router) {
	router.Post("/user", cuser.Register)
}
